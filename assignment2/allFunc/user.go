package allFunc

import (
	"assignment2/models"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	UserController struct {
		session *mgo.Session
	}
)

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

func (uc UserController) GetLocations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	objectId := bson.ObjectIdHex(id)
	location := structure.Location{}
	if err := uc.session.DB("cmpe273_project").C("hello").FindId(objectId).One(&location); err != nil {
		w.WriteHeader(404)
		return
	}

	result, _ := json.Marshal(location)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", result)
}

func (userCon UserController) CreateLocations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	location := structure.Location{}

	json.NewDecoder(r.Body).Decode(&location)
	str := getURL(location.Address, location.City, location.State)
	getLocation(&location, str)

	location.Id = bson.NewObjectId()
	userCon.session.DB("cmpe273_project").C("hello").Insert(location)

	result, _ := json.Marshal(location)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", result)
}

func (uc UserController) UpdateLocations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)
	l := structure.Location{}
	json.NewDecoder(r.Body).Decode(&l)
	str := getURL(l.Address, l.City, l.State)
	getLocation(&l, str)
	l.Id = oid
	if err := uc.session.DB("cmpe273_project").C("hello").Update(bson.M{"_id": l.Id}, bson.M{"$set": bson.M{"address": l.Address, "city": l.City, "state": l.State, "zip": l.Zip, "coordinate.lat": l.Coordinate.Lat, "coordinate.lng": l.Coordinate.Lng}}); err != nil {
		w.WriteHeader(404)
		return
	}
	if err := uc.session.DB("cmpe273_project").C("hello").FindId(oid).One(&l); err != nil {
		w.WriteHeader(404)
		return
	}
	lj, _ := json.Marshal(l)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", lj)
}

func (uc UserController) RemoveLocations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)

	if err := uc.session.DB("cmpe273_project").C("hello").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)

}

func getURL(address string, city string, state string) string {
	addStr := address
	cityStr := city
	stateStr := state
	add := strings.Split(addStr, " ")
	var res string
	for i := 0; i < len(add); i++ {
		if i == len(add)-1 {
			res = res + add[i] + ","
		} else {
			res = res + add[i] + "+"
		}
	}
	c := strings.Split(cityStr, " ")
	for i := 0; i < len(c); i++ {
		if i == len(c)-1 {
			res = res + "+" + c[i] + ","
		} else {
			res = res + "+" + c[i]
		}
	}
	res = res + "+" + stateStr
	return res
}

func getLocation(l *structure.Location, str string) {
	urlPath := "http://maps.google.com/maps/api/geocode/json?address="
	urlPath += str
	urlPath += "&sensor=false"
	res, err := http.Get(urlPath)
	if err != nil {
		fmt.Println("GetLocation: http.Get", err)
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("GetLocation: ioutil.ReadAll", err)
		panic(err)
	}

	mp := make(map[string]interface{})
	err = json.Unmarshal(body, &mp)
	if err != nil {
		fmt.Println("GetLocation: json.Unmarshal", err)
		panic(err)
	}
	temp := mp["results"].(interface{})
	next := temp.([]interface{})
	geometry := next[0].(map[string]interface{})
	geometry = geometry["geometry"].(map[string]interface{})
	location := geometry["location"].(map[string]interface{})
	lat := location["lat"].(float64)
	lng := location["lng"].(float64)
	l.Coordinate.Lat = lat
	l.Coordinate.Lng = lng
}
