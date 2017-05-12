package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type ClassifierTypeSelector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SurveySummary struct {
	ID     string `json:"id"`
	Survey string `json:"survey"`
}

type Survey struct {
	ID     string `json:"id"`
	Survey string `json:"survey"`
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
	router.GET("/surveys/:surveyid", getSurveyEndpoint)
	router.GET("/surveys/:surveyid/classifiertypeselectors", listClassifierTypeSelectorsEndpoint)

	router.Run(port)
}

// GET /classifiertypes
func listClassifierTypeSelectorsEndpoint(context *gin.Context) {
	classifierTypeSelectors := getClassifierTypeSelectors(context.Param("surveyid"))

	if len(classifierTypeSelectors) == 0 {
		context.AbortWithStatus(http.StatusNoContent)
	} else {
		context.JSON(http.StatusOK, classifierTypeSelectors)
	}
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

// GET /surveys/{surveyid}
func getSurveyEndpoint(context *gin.Context) {
	survey := getSurvey(context.Param("surveyid"))
	context.JSON(http.StatusOK, survey)
}

func getClassifierTypeSelectors(surveyID string) []ClassifierTypeSelector {
	rows, err := db.Query("SELECT classifiertypeselector.id, classifiertypeselector FROM survey.classifiertypeselector INNER JOIN survey.survey ON classifiertypeselector.surveyid = survey.surveyid WHERE survey.id = $1 ORDER BY classifiertypeselector", surveyID)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var classifierTypeSelectors []ClassifierTypeSelector
	var classifierTypeSelector ClassifierTypeSelector

	for rows.Next() {
		if err := rows.Scan(&classifierTypeSelector.ID, &classifierTypeSelector.Name); err != nil {
			log.Fatal(err)
		}

		classifierTypeSelectors = append(classifierTypeSelectors, classifierTypeSelector)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return classifierTypeSelectors
}

func getSurvey(surveyID string) Survey {
	rows, err := db.Query("SELECT id, survey from survey.survey WHERE id = $1", surveyID)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var survey Survey

	for rows.Next() {
		if err := rows.Scan(&survey.ID, &survey.Survey); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return survey
}

func getSurveys() []SurveySummary {
	rows, err := db.Query("SELECT id, survey FROM survey.survey ORDER BY survey ASC")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var surveySummaries []SurveySummary

	for rows.Next() {
		var surveySummary SurveySummary

		if err := rows.Scan(&surveySummary.ID, &surveySummary.Survey); err != nil {
			log.Fatal(err)
		}

		surveySummaries = append(surveySummaries, surveySummary)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return surveySummaries
}
