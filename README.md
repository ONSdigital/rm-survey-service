# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/) and has the following responsibilities:

* Providing a list of available surveys
* Providing the details for a survey (including classifier types)

## Prerequisites
* Install the [Go PostgreSQL driver]() using `go get github.com/lib/pq`
* Install the [Gorilla Mux URL router](http://www.gorillatoolkit.org/pkg/mux) using `go get github.com/gorilla/mux`

## Building
### Docker Image
To build the Docker image, from the project root run:

```
docker build -t surveysvc .
```

## Running
From `$GOPATH`, use `go run src/github.com/onsdigital/rm-survey-service/survey-api/main.go &` to start the Survey service in the background. The following environment variables may be overridden:

| Environment Variable | Purpose                               | Default Value                                                   |
| :------------------- | :------------------------------------ | :-------------------------------------------------------------- |
| DATABASE_URL         | PostgreSQL database connection string | postgres://postgres:password@localhost/postgres?sslmode=disable |
| PORT                 | HTTP listener port                    | :8080                                                           |

### Docker Image and PostgreSQL
To start the Docker image, run:

```
docker-compose up -d
```

Initial data can be loaded into the PostgreSQL database by starting the Docker image and connecting, then loading the SQL files in the `sql` directory:
```
docker exec -it postgres /bin/sh
psql -U postgres -d postgres -f ./sql/groundzero.sql
psql -U postgres -d postgres -f ./sql/survey_foundation_schema.sql
psql -U postgres -d postgres -f ./sql/seed_data.sql
psql postgres://postgres:password@localhost/postgres?sslmode=disable
```

## API Examples
### List Surveys

* Running the command `curl http://localhost:8080/surveys` should return an HTTP 200 status code with the JSON response:

```json
["BRES"]
```

### Get Survey

* Running the command `curl http://localhost:8080/surveys/bres` should return an HTTP 200 status code with the JSON response:

```json
{
  "survey": "BRES",
  "classifierTypes": [
    "COLLECTION_EXERCISE",
    "LEGAL_BASIS",
    "RU_REF",
    "SIC"
  ]
}
```

## Testing
To follow once I've worked out how to write unit tests in Go :-)

## Copyright
Copyright (C) 2017 Crown Copyright (Office for National Statistics)
