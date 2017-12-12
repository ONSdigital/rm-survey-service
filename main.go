package main

import (
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

const (
	realm       = "sdc"
	serviceName = "surveysvc"
	timeFormat  = "2006-01-02T15:04:05Z0700"
)

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

	r.HandleFunc("/info", surveys.Info).Methods("GET")
	r.HandleFunc("/surveys", use(surveys.AllSurveys, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}", use(surveys.GetSurveyByIdHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{shortName}", use(surveys.GetSurveyByShortName, basicAuth)).
		Methods("GET")
	r.HandleFunc("/surveys/ref/{ref}", use(surveys.GetSurveyByReference, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors", use(surveys.AllClassifierTypeSelectors, BasicAuth)).
		Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors/{classifierTypeSelectorId}",
		use(surveys.GetClassifierTypeSelector, basicAuth)).
		Method("GET")

	http.Handle("/", r)

	// CompressHandler gzips responses where possible
	compressHandler := handlers.CompressHandler()
	log.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), compressHandler))
}

func use(h http.HandlerFunc, m middleware, f ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
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
		LogInfo("Found Cloud Foundry environment")

		if err == nil {
			uri, found := postgresServer[0].CredentialString("uri")

			if found {
				dataSource = uri
			} else {
				message := "Unable to retrieve PostgreSQL URI from Cloud Foundry environment"
				LogInfo(message)
				panic(message)
			}
		}
	} else {
		LogInfo("No Cloud Foundry environment")

		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}

		if v := os.Getenv("DATABASE_URL"); len(v) > 0 {
			dataSource = v
		}
	}

	return dataSource, port
}

func LogError(message string, err error) {
	logger.Error(message,
		zap.String("service", serviceName),
		zap.String("event", "error"),
		zap.String("data", err.Error()),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}

func LogInfo(message string) {
	logger.Info(message,
		zap.String("service", serviceName),
		zap.String("event", "service startup"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}
