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
	dataSource, port, migrationSource := configureEnvironment()
	db, err := models.InitDB(dataSource, migrationSource)

	if err != nil {
		logger.Fatal(fmt.Sprintf(`event="Failed to start" error="unable to initialise database" error_message=%s`, err.Error()))
	}

	// Set up the signal handler to watch for SIGTERM and SIGINT signals so we
	// can at least attempt to gracefully shut down before the PaaS/docker etc
	// running us unceremoneously kills us with a SIGKILL.

	api, err := models.NewAPI(db)
	if err != nil {
		logger.Fatal(fmt.Sprintf(`event="Failed to start" error="unable to initialise API model" error_message=%s`, err.Error()))
	}

	// Webserver - strictslash set to true to match trailing slashes to routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/info", api.Info).Methods("GET")
	r.HandleFunc("/surveys", use(api.AllSurveys, basicAuth)).Methods("GET")
	r.HandleFunc("/legal-bases", use(api.AllLegalBases, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}", use(api.GetSurvey, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/shortname/{shortName}", use(api.GetSurveyByShortName, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/ref/{ref}", use(api.PutSurveyDetails, basicAuth)).Methods("PUT")
	r.HandleFunc("/surveys", use(api.PostSurveyDetails, basicAuth)).Methods("POST")
	r.HandleFunc("/surveys/ref/{ref}", use(api.GetSurveyByReference, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors", use(api.AllClassifierTypeSelectors, basicAuth)).Methods("GET")
	r.HandleFunc("/surveys/{surveyId}/classifiertypeselectors/{classifierTypeSelectorId}", use(api.GetClassifierTypeSelectorByID, basicAuth)).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handlers.CompressHandler(r),
	}
	log.Fatalf(`event="Stopped" error="%v"`, srv.ListenAndServe())
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
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if pair[0] != os.Getenv("security_user_name") || pair[1] != os.Getenv("security_user_password") {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func configureEnvironment() (dataSource, port string, migrationSource string) {
	dataSource = "postgres://postgres:password@localhost/postgres?sslmode=disable"
	port = "8080"
	migrationSource = "file:///db-migrations"
    if v := os.Getenv("MIGRATION_SOURCE"); len(v) > 0 {
        migrationSource = v
    }

	appEnv, err := cfenv.Current()

	if err != nil {
		LogInfo("No Cloud Foundry environment")

		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}

		if v := os.Getenv("DATABASE_URL"); len(v) > 0 {
			dataSource = v
		}

		return dataSource, port, migrationSource
	}

	ps := appEnv.Port
	port = strconv.FormatInt(int64(ps), 10)
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

	return dataSource, port, migrationSource
}

//LogError log out error messages
func LogError(message string, err error) {
	logger.Error(message,
		zap.String("service", serviceName),
		zap.String("event", "error"),
		zap.String("data", err.Error()),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}

//LogInfo log out info log messages
func LogInfo(message string) {
	logger.Info(message,
		zap.String("service", serviceName),
		zap.String("event", "service startup"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}
