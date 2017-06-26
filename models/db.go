package models

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const serviceName = "surveysvc"
const timeFormat = "2006-01-02T15:04:05Z0700"

var db *sql.DB
var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func InitDB(dataSource string) {
	const DriverName = "postgres"
	var err error
	db, err = sql.Open(DriverName, dataSource)
	if err != nil {
		logger.Error("Error opening data source",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))
	}

	if err = db.Ping(); err != nil {
		logger.Error("Error establishing connection to data source",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))
	}

	if !schemaExists() {
		bootstrapSchema()
	}
}

func bootstrapSchema() {
	exe, _ := os.Executable()
	exePath := path.Dir(exe)
	file, err := ioutil.ReadFile(exePath + "/sql/bootstrap.sql")
	if err != nil {
		logger.Error("Error reading bootstrap schema SQL file",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))
	}

	logger.Info("Creating and populating database schema",
		zap.String("service", serviceName),
		zap.String("event", "database bootstrap"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))

	statements := strings.Split(string(file), ";")
	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			logger.Error("Error executing bootstrap SQL statement",
				zap.String("service", serviceName),
				zap.String("event", "error"),
				zap.String("data", err.Error()),
				zap.String("created", time.Now().UTC().Format(timeFormat)))
		}
	}
}

func schemaExists() bool {
	var schemaName string
	err := db.QueryRow("SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'survey'").Scan(&schemaName)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("Database schema doesn't exist",
				zap.String("service", serviceName),
				zap.String("event", "database bootstrap"),
				zap.String("created", time.Now().UTC().Format(timeFormat)))

			return false
		}

		logger.Error("Error executing schema exists check SQL statement",
			zap.String("service", serviceName),
			zap.String("event", "error"),
			zap.String("data", err.Error()),
			zap.String("created", time.Now().UTC().Format(timeFormat)))
	}

	logger.Info("Database schema exists",
		zap.String("service", serviceName),
		zap.String("event", "database bootstrap"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))

	return true
}
