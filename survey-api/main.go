package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const contentTypeHeader string = "Content-Type"
const contentType string = "application/json"

var db *sql.DB
var err error

func main() {
	port := ":8080"
	dataSource := "postgres://postgres:password@localhost/postgres?sslmode=disable"

	if v := os.Getenv("PORT"); len(v) > 0 {
		port = v
	}

	if v := os.Getenv("DATABASE_URL"); len(v) > 0 {
		dataSource = v
	}

	db, err = sql.Open("postgres", dataSource)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
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
	rows, err := db.Query("SELECT survey FROM survey.survey ORDER BY survey ASC")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var surveys []string

	for rows.Next() {
		var survey string

		if err := rows.Scan(&survey); err != nil {
			log.Fatal(err)
		}

		surveys = append(surveys, survey)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set(contentTypeHeader, contentType)
	json.NewEncoder(w).Encode(surveys)
}

// GET /surveys/{survey}/classifiertypes
func classifierTypesHandler(w http.ResponseWriter, req *http.Request) {
	survey := mux.Vars(req)["survey"]
	log.Printf("Getting the list of classifier types for survey '%s'", survey)
}
