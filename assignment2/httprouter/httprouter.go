package main

import (
	"assignment2/allFunc"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
)

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/cmpe273_project")
	if err != nil {
		fmt.Println("Can't connect to mongo, error %v\n", err)
	}
	return s
}

func main() {
	mux := httprouter.New()
	helper := allFunce.NewUserController(getSession())
	mux.GET("/locations/:id", helper.GetLocations)
	mux.POST("/locations/", helper.CreateLocations)
	mux.DELETE("/locations/:id", helper.RemoveLocations)
	mux.PUT("/locations/:id", helper.UpdateLocations)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
