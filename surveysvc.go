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

		logger.Info("Found Cloud Foundry environment",
			zap.String("service", serviceName),
			zap.String("event", "service startup"),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		if err == nil {
			uri, found := postgresServer[0].CredentialString("uri")

			if found {
				dataSource = uri
			}
		}
	} else {
		logger.Info("No Cloud Foundry environment",
			zap.String("service", serviceName),
			zap.String("event", "service startup"),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

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

	echo.GET("/surveys", allSurveys)
	echo.GET("/surveys/:surveyid", getSurvey)
	echo.GET("/surveys/name/:name", getSurveyByName)
	echo.GET("/surveys/ref/:ref", getSurveyByReference)
	echo.GET("/surveys/:surveyid/classifiertypeselectors", allClassifierTypeSelectors)
	echo.GET("/surveys/:surveyid/classifiertypeselectors/:classifiertypeselectorid", getClassifierTypeSelector)

	logger.Info("Survey service started on port "+port,
		zap.String("service", serviceName),
		zap.String("event", "service started"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))

	echo.Start(port)
}

func allSurveys(context echo.Context) error {
	surveys, err := models.AllSurveys()
	if err != nil {
		logger.Error("Error getting list of surveys",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

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
		}

		logger.Error("Error getting survey '"+surveyID+"'",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		return context.JSON(http.StatusInternalServerError, err.Error())

	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByName(context echo.Context) error {
	name := context.Param("name")
	survey, err := models.GetSurveyByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		}

		logger.Error("Error getting survey '"+name+"'",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByReference(context echo.Context) error {
	reference := context.Param("ref")
	survey, err := models.GetSurveyByReference(reference)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		}

		logger.Error("Error getting survey '"+reference+"'",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

func allClassifierTypeSelectors(context echo.Context) error {
	surveyID := context.Param("surveyid")
	classifierTypeSelectors, err := models.AllClassifierTypeSelectors(surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return context.JSON(http.StatusNotFound, "Survey not found")
		}

		logger.Error("Error getting list of classifier type selectors for survey '"+surveyID+"'",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

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
			return context.JSON(http.StatusNotFound, "Survey or classifier type selector not found")
		}

		logger.Error("Error getting classifier type selector '"+classifierTypeSelectorID+"' for survey '"+surveyID+"'",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))

		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, classifierTypeSelector)
}
