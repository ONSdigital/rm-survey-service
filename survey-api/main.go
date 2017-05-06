package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const contentTypeHeader string = "Content-Type"
const contentType string = "application/json"

var db *sql.DB
var err error

type classifierTypes struct {
	Survey          string   `json:"survey"`
	ClassifierTypes []string `json:"classifierTypes"`
}

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
	var responseJSON []string

	for rows.Next() {
		var survey string

		if err := rows.Scan(&survey); err != nil {
			log.Fatal(err)
		}

		responseJSON = append(responseJSON, survey)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set(contentTypeHeader, contentType)
	json.NewEncoder(w).Encode(responseJSON)
}

// GET /surveys/{survey}/classifiertypes
func classifierTypesHandler(w http.ResponseWriter, req *http.Request) {
	survey := mux.Vars(req)["survey"]
	log.Printf("Getting the list of classifier types for survey '%s'", survey)

	classifierTypes := getClassifierTypes(strings.ToUpper(survey))
	b, err := json.Marshal(&classifierTypes)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set(contentTypeHeader, contentType)
	w.Write(b)
}

func getClassifierTypes(survey string) classifierTypes {
	rows, err := db.Query("SELECT classifiertype FROM survey.classifiertype INNER JOIN survey.survey ON classifiertype.surveyid = survey.surveyid WHERE survey= $1 ORDER BY classifiertype ASC", survey)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var responseJSON classifierTypes
	var classifierTypes []string

	for rows.Next() {
		var classifierType string

		if err := rows.Scan(&classifierType); err != nil {
			log.Fatal(err)
		}

		classifierTypes = append(classifierTypes, classifierType)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	responseJSON.Survey = survey
	responseJSON.ClassifierTypes = classifierTypes

	return responseJSON
}
