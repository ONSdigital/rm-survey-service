package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
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

// InitDB opens dataSource and bootstraps the database schema if it doesn't already exist.
func InitDB(dataSource string) {
	const DriverName = "postgres"
	var err error
	db, err = sql.Open(DriverName, dataSource)
	if err != nil {
		logError("Error opening data source", err)
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
		logError(fmt.Sprintf("Error reading '%s/sql/bootstrap.sql' file", exePath), err)
	}

	logInfo("Creating and populating database schema")
	statements := strings.Split(string(file), ";")

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			logError(fmt.Sprintf("Error executing '%s/sql/bootstrap.sql file", exePath), err)
		}
	}
}

func schemaExists() bool {
	var schemaName string
	err := db.QueryRow("SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'survey'").Scan(&schemaName)
	if err != nil {
		if err == sql.ErrNoRows {
			logInfo("Database schema doesn't exist")
			return false
		}

		logError("Error executing schema exists check SQL statement", err)
	}

	logInfo("Database schema exists")
	return true
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
		zap.String("event", "database bootstrap"),
		zap.String("created", time.Now().UTC().Format(timeFormat)))
}
