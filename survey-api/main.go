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
	router.GET("/surveys/:surveyid/classifiertypeselectors/:classifiertypeselectorid", getClassifierTypeSelector)

	router.Run(port)
}

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

func getSurvey(context *gin.Context) {
	surveyID := context.Param("surveyid")
	survey, err := models.GetSurvey(surveyID)
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

func allClassifierTypeSelectors(context *gin.Context) {
	surveyID := context.Param("surveyid")
	classifierTypeSelectors, err := models.AllClassifierTypeSelectors(surveyID)
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

func getClassifierTypeSelector(context *gin.Context) {
	surveyID := context.Param("surveyid")
	classifierTypeSelectorID := context.Param("classifiertypeselectorid")
	classifierTypeSelector, err := models.GetClassifierTypeSelector(surveyID, classifierTypeSelectorID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, "Survey or classifier type selector not found")
		} else {
			log.Println(err)
			context.JSON(http.StatusInternalServerError, err.Error())
		}

		return
	}

	context.JSON(http.StatusOK, classifierTypeSelector)
}
