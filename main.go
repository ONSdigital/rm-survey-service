package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/ONSdigital/rm-survey-service/models"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/handlers"
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

	api, err := models.NewAPI(db)
	if err != nil {
		logger.Fatal(fmt.Sprintf(`event="Failed to start" error="unable to initialise API model" error_message=%s`, err.Error()))
	}

	r := models.Router(api)

	http.Handle("/", r)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handlers.CompressHandler(r),
	}
	log.Fatalf(`event="Stopped" error="%v"`, srv.ListenAndServe())
}

func configureEnvironment() (dataSource, port string) {
	dataSource = "postgres://postgres:password@localhost/postgres?sslmode=disable"
	port = "8080"
	appEnv, err := cfenv.Current()

	if err != nil {
		LogInfo("No Cloud Foundry environment")

		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}

		if v := os.Getenv("DATABASE_URL"); len(v) > 0 {
			dataSource = v
		}

		return dataSource, port
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

	return dataSource, port
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
