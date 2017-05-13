package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/onsdigital/rm-survey-service/survey-api/models"
)

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

	models.InitDB(dataSource)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/surveys", allSurveys)
	router.GET("/surveys/:surveyid", getSurvey)
	router.GET("/surveys/:surveyid/classifiertypeselectors", allClassifierTypeSelectors)

	router.Run(port)
}

// GET /surveys
func allSurveys(context *gin.Context) {
	surveys, err := models.AllSurveys()
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(surveys) == 0 {
		context.AbortWithStatus(http.StatusNoContent)
	} else {
		context.JSON(http.StatusOK, surveys)
	}
}

// GET /surveys/{surveyid}
func getSurvey(context *gin.Context) {
	survey, err := models.GetSurvey(context.Param("surveyid"))
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, "Survey not found")
		} else {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
		}

		return
	}

	context.JSON(http.StatusOK, survey)
}

// GET /surveys/{surveyid}/classifiertypeselectors
func allClassifierTypeSelectors(context *gin.Context) {
	classifierTypeSelectors, err := models.AllClassifierTypeSelectors(context.Param("surveyid"))
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, "Survey not found")
		} else {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
		}

		return
	}

	if len(classifierTypeSelectors) == 0 {
		context.AbortWithStatus(http.StatusNoContent)
	} else {
		context.JSON(http.StatusOK, classifierTypeSelectors)
	}
}
