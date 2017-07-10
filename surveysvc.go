package main // import "github.com/onsdigital/rm-survey-service"

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onsdigital/rm-survey-service/models"
)

const serviceName = "surveysvc"
const timeFormat = "2006-01-02T15:04:05Z0700"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func main() {
	port := ":8080"
	dataSource := "postgres://postgres:password@localhost/postgres?sslmode=disable"
	appEnv, err := cfenv.Current()

	if err == nil {
		ps := appEnv.Port
		port = ":" + strconv.FormatInt(int64(ps), 10)
		postgresServer, err := appEnv.Services.WithTag("postgresql")
		logInfo("Found Cloud Foundry environment")

		if err == nil {
			uri, found := postgresServer[0].CredentialString("uri")

			if found {
				dataSource = uri
			}
		}
	} else {
		logInfo("No Cloud Foundry environment")

		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}

		if v := os.Getenv("DATABASE_URL"); len(v) > 0 {
			dataSource = v
		}
	}

	models.InitDB(dataSource)

	echo := echo.New()
	echo.Use(middleware.Gzip())

	echo.GET("/info", info)
	echo.GET("/surveys", allSurveys)
	echo.GET("/surveys/:surveyid", getSurvey)
	echo.GET("/surveys/name/:name", getSurveyByName)
	echo.GET("/surveys/ref/:ref", getSurveyByReference)
	echo.GET("/surveys/:surveyid/classifiertypeselectors", allClassifierTypeSelectors)
	echo.GET("/surveys/:surveyid/classifiertypeselectors/:classifiertypeselectorid", getClassifierTypeSelector)

	logInfo("Survey service started on port " + port)
	echo.Start(port)
}

func info(context echo.Context) error {
	return context.JSON(http.StatusOK, models.NewVersion())
}

func allSurveys(context echo.Context) error {
	surveys, err := models.AllSurveys()
	if err != nil {
		logError("Error getting list of surveys", err)
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
			re := models.NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())

	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByName(context echo.Context) error {
	name := context.Param("name")
	survey, err := models.GetSurveyByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			re := models.NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting survey '"+name+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByReference(context echo.Context) error {
	reference := context.Param("ref")
	survey, err := models.GetSurveyByReference(reference)
	if err != nil {
		if err == sql.ErrNoRows {
			re := models.NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting survey '"+reference+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

func allClassifierTypeSelectors(context echo.Context) error {
	surveyID := context.Param("surveyid")
	classifierTypeSelectors, err := models.AllClassifierTypeSelectors(surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			re := models.NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting list of classifier type selectors for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
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
			re := models.NewRESTError("404", "Survey or classifier type selector not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, classifierTypeSelector)
}

func logError(message string, err error) {
	logger.Error(message,
		zap.String("service", serviceName),
		zap.String("event", "error"),
		zap.String("data", err.Error()),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}

func logInfo(message string) {
	logger.Info(message,
		zap.String("service", serviceName),
		zap.String("event", "service startup"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}
