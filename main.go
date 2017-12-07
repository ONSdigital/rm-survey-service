package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/onsdigital/rm-survey-service/models"
)

const realm = "sdc"
const serviceName = "surveysvc"
const timeFormat = "2006-01-02T15:04:05Z0700"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func main() {
	dataSource, port := configureEnvironment()
	models.InitDB(dataSource)

	// Set up the signal handler to watch for SIGTERM and SIGINT signals so we
	// can at least attempt to gracefully shut down before the PaaS/docker etc
	// running us unceremoneously kills us with a SIGKILL.
	cancelSigWatch := signals.HandleFunc(
		func(sig os.Signal) {
			log.Printf(`event="Shutting down" signal="%s"`, sig.String())
			// If any clean-up of attached services is needed, do it here
			log.Print(`event="Exiting"`)
			os.Exit(0)
		},
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancelSigWatch()

	// Webserver - strictslash set to true to match trailing slashes to routes
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/info", infoHandler).Methods("GET")
	r.HandleFunc("/surveys", use(getAllSurveysHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyid}", use(getSurveyByIdHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{shortname}", use(getSurveyByShortNameHandler, basicAuth)).
		Methods("GET")
	r.HandleFunc("/surveys/ref/{ref}", use(getSurveyByRefHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyid}/classifiertypeselectors", use(getSurveyClassifierHandler, BasicAuth)).
		Methods("GET")
	r.HandleFunc("/surveys/{surveyid}/classifiertypeselectors/{classifiertypeselectorid}",
		use(getSurveyClassifierById, basicAuth)).
		Method("GET")

	http.Handle("/", r)

	// CompressHandler gzips responses where possible
	compressHandler := handlers.CompressHandler()
	log.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), compressHandler))
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	// Taken from https://gist.github.com/elithrar/9146306
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != os.Getenv("username") || pair[1] != os.Getenv("password") {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func configureEnvironment() (dataSource, port string) {
	dataSource = "postgres://postgres:password@localhost/postgres?sslmode=disable"
	port = ":8080"
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
			} else {
				message := "Unable to retrieve PostgreSQL URI from Cloud Foundry environment"
				logInfo(message)
				panic(message)
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

	return dataSource, port
}

func info(w http.ResponseWriter, r *http.Request) error {
	return context.JSON(http.StatusOK, models.NewVersion())
}

func allSurveys() error {
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

func getSurvey() error {
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

func getSurveyByShortName() error {
	shortName := context.Param("shortname")
	survey, err := models.GetSurveyByShortName(shortName)
	if err != nil {
		if err == sql.ErrNoRows {
			re := models.NewRESTError("404", "Survey not found")
			return context.JSON(http.StatusNotFound, re)
		}

		logError("Error getting survey '"+shortName+"'", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, survey)
}

func getSurveyByReference() error {
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

func allClassifierTypeSelectors() error {
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

func getClassifierTypeSelector() error {
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
