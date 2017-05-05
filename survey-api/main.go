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
	subRouter.HandleFunc("/{survey}/classifiertypes", classifierTypesHandler)

	log.Printf("Survey service listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func classifierTypesHandler(w http.ResponseWriter, req *http.Request) {

}
