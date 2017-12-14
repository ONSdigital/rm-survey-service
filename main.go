package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ONSdigital/rm-survey-service/models"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	db := models.InitDB(dataSource)

	// Set up the signal handler to watch for SIGTERM and SIGINT signals so we
	// can at least attempt to gracefully shut down before the PaaS/docker etc
	// running us unceremoneously kills us with a SIGKILL.

	api, _ := models.NewAPI(db)

	// Webserver - strictslash set to true to match trailing slashes to routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/info", api.Info).Methods("GET")
	r.HandleFunc("/surveys", use(api.AllSurveys, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}", use(api.GetSurvey, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{surveyId}", use(api.GetSurvey, basicAuth)).Methods("GET")

	http.Handle("/", r)

	// CompressHandler gzips responses where possible
	compressHandler := handlers.CompressHandler(r)
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