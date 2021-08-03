package models

import (
	"database/sql"
	"time"

	"github.com/blendle/zapdriver"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"go.uber.org/zap"
)

const serviceName = "surveysvc"
const timeFormat = "2006-01-02T15:04:05Z0700"

var db *sql.DB
var logger *zap.Logger

func init() {
	logger, _ = zapdriver.NewProduction()
	defer logger.Sync()
}

// InitDB opens dataSource and bootstraps the database schema if it doesn't already exist.
func InitDB(dataSource string, migrationSource string, maxIdleConn int, connMaxLifetime int) (*sql.DB, error) {
	const DriverName = "postgres"
	var err error
	db, err = sql.Open(DriverName, dataSource)

	if err != nil {
		logError("Error opening data source", err)
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	err = testDBConnection(db)

	if err != nil {
		logError("Error pinging data source", err)
		return nil, err
	}

	err = bootstrapSchema(db, migrationSource)

	if err == migrate.ErrNoChange {
		logInfo("Database schema unchanged")
		err = nil
	}

	return db, err
}

func testDBConnection(db *sql.DB) error {
	// Keep attempting to ping the database, increasing the time between each attempt.
	// See https://medium.com/@kelseyhightower/12-fractured-apps-1080c73d481c
	maxAttempts := 20
	var err error

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			break
		}

		logError("Error pinging data source", err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	return err
}

func bootstrapSchema(db *sql.DB, migrationSource string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "survey_schema_migrations"})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
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
