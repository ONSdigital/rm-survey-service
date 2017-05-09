package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type SurveySummary struct {
	UUID   string `json:"id"`
	Survey string `json:"survey"`
}

type survey struct {
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

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/surveys", listSurveysEndpoint)
	router.GET("/surveys/:survey", getSurveyEndpoint)

	router.Run(port)
}

// GET /surveys
func listSurveysEndpoint(context *gin.Context) {
	surveys := getSurveys()

	if len(surveys) == 0 {
		context.AbortWithStatus(http.StatusNoContent)
	} else {
		context.JSON(http.StatusOK, surveys)
	}
}

// GET /surveys/{survey}
func getSurveyEndpoint(context *gin.Context) {
	survey := getSurvey(strings.ToUpper(context.Param("survey")))
	context.JSON(http.StatusOK, survey)
}

func getSurvey(surveyName string) survey {
	rows, err := db.Query("SELECT classifiertype FROM survey.classifiertype INNER JOIN survey.survey ON classifiertype.surveyid = survey.surveyid WHERE survey= $1 ORDER BY classifiertype ASC", surveyName)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var survey survey
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

	survey.Survey = surveyName
	survey.ClassifierTypes = classifierTypes

	return survey
}

func getSurveys() []SurveySummary {
	rows, err := db.Query("SELECT uuid, survey FROM survey.survey ORDER BY survey ASC")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var surveySummaries []SurveySummary

	for rows.Next() {
		var surveySummary SurveySummary

		if err := rows.Scan(&surveySummary.UUID, &surveySummary.Survey); err != nil {
			log.Fatal(err)
		}

		surveySummaries = append(surveySummaries, surveySummary)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return surveySummaries
}
