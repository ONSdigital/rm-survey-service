package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := ":8080"

	if v := os.Getenv("SURVEY_SERVICE_PORT"); len(v) > 0 {
		port = v
	}

	// If there's a trailing slash, redirect to the non-trailing slash URL.
	router := mux.NewRouter().StrictSlash(true)
	subRouter := router.PathPrefix("/surveys").Subrouter()
	subRouter.HandleFunc("", surveysHandler)
	subRouter.HandleFunc("/{survey}/classifiertypes", classifierTypesHandler)

	log.Printf("Survey service listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// GET /surveys
func surveysHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Getting the list of surveys")
}

// GET /surveys/{survey}/classifiertypes
func classifierTypesHandler(w http.ResponseWriter, req *http.Request) {
	survey := mux.Vars(req)["survey"]
	log.Printf("Getting the list of classifier types for survey '%s'", survey)
}
