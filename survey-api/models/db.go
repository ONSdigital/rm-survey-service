package models

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(adminDataSource string, dataSource string) {
	const DriverName = "postgres"
	var err error
	db, err = sql.Open(DriverName, adminDataSource)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	if !schemaExists() {
		bootstrapSchema()
	}

	db, err = sql.Open(DriverName, dataSource)
	if err != nil {
		log.Panic(err)
	}
}

func bootstrapSchema() {
	file, err := ioutil.ReadFile("sql/bootstrap.sql")
	if err != nil {
		log.Println(err)
	}

	statements := strings.Split(string(file), ";")
	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			log.Println(err)
		}
	}
}

func schemaExists() bool {
	var schemaName string
	err := db.QueryRow("SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'survey'").Scan(&schemaName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Println(err)
	}

	return true
}
