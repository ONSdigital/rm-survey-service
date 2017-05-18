package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	echo := echo.New()
	echo.Use(middleware.Gzip())

	echo.GET("/surveys", allSurveys)
	echo.GET("/surveys/:surveyid", getSurvey)
	echo.GET("/surveys/name/:name", getSurveyByName)
	echo.GET("/surveys/:surveyid/classifiertypeselectors", allClassifierTypeSelectors)
	echo.GET("/surveys/:surveyid/classifiertypeselectors/:classifiertypeselectorid", getClassifierTypeSelector)

	echo.Logger.Fatal(echo.Start(port))
}

func allSurveys(context echo.Context) error {
	surveys, err := models.AllSurveys()
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(surveys) == 0 {
		return context.String(http.StatusNoContent, "")
	}

	return context.JSON(http.StatusOK, surveys)
}

func getSurvey(context echo.Context) error {
	surveyID := context.Param("surveyid")
	survey, err := models.GetSurvey(surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		} else {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByName(context echo.Context) error {
	name := context.Param("name")
	survey, err := models.GetSurveyByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		} else {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return context.JSON(http.StatusOK, survey)
}

func allClassifierTypeSelectors(context echo.Context) error {
	surveyID := context.Param("surveyid")
	classifierTypeSelectors, err := models.AllClassifierTypeSelectors(surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		} else {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if len(classifierTypeSelectors) == 0 {
		return context.String(http.StatusNoContent, "")
	}

	return context.JSON(http.StatusOK, classifierTypeSelectors)
}

func getClassifierTypeSelector(context echo.Context) error {
	surveyID := context.Param("surveyid")
	classifierTypeSelectorID := context.Param("classifiertypeselectorid")
	classifierTypeSelector, err := models.GetClassifierTypeSelector(surveyID, classifierTypeSelectorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey or classifier type selector not found")
		} else {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return context.JSON(http.StatusOK, classifierTypeSelector)
}
