package models

import (
	"database/sql"
	"fmt"
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
func InitDB(dataSource string) *sql.DB {
	const DriverName = "postgres"
	var err error
	db, err = sql.Open(DriverName, dataSource)
	if err != nil {
		logError("Error opening data source", err)
	}

	// Keep attempting to ping the database, increasing the time between each attempt.
	// See https://medium.com/@kelseyhightower/12-fractured-apps-1080c73d481c
	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			break
		}

		logError("Error pinging data source", err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if err != nil {
		logError("Error pinging data source", err)
	}
	if !schemaExists() {
		bootstrapSchema()
	}
	return db
}

func bootstrapSchema() {
	logInfo("Creating and populating database schema")
	sql := bootstrapSQL()

	for _, query := range sql {
		_, err := db.Exec(query)
		if err != nil {
			logError(fmt.Sprintf("Error executing bootstrap statement: '%s'", query), err)
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
