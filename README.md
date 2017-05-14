# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/) and has the following responsibilities:

* Providing a list of known surveys
* Providing the details for a specified survey

## Prerequisites
* Install the [Go PostgreSQL driver](https://github.com/lib/pq) using `go get github.com/lib/pq`
* Install the [Gin HTTP web framework](https://gin-gonic.github.io/gin/) using `go get gopkg.in/gin-gonic/gin.v1`
* Install the [Gin GZIP middleware](https://github.com/gin-contrib/gzip) using `go get github.com/gin-contrib/gzip`

## Building
### Make
A Makefile is provided for compiling the code using `go build`:

```
make
```

The compiled executable is placed within the `build` directory tree.

### Docker Image
To build the Docker image, from the project root run:

```
docker build -t surveysvc .
```

## Running
From `$GOPATH`, use `go run src/github.com/onsdigital/rm-survey-service/survey-api/main.go &` to start the Survey service in the background. Or compile the service first using `make` and execute the binary in the background using `./surveysvc &` from within the `bin` directory within the `build` directory tree.

The following environment variables may be overridden:

| Environment Variable | Purpose                               | Default Value                                                   |
| :------------------- | :------------------------------------ | :-------------------------------------------------------------- |
| DATABASE_URL         | PostgreSQL database connection string | postgres://postgres:password@localhost/postgres?sslmode=disable |
| GIN_MODE             | Gin debug/release/test mode           | debug                                                           |
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
```

## API Examples
### List Surveys

* Running the command `curl http://localhost:8080/surveys` should return an HTTP 200 status code with the JSON response:

```json
[{
    "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
    "survey": "BRES"
}]
```

### Get Survey

* Running the command `curl http://localhost:8080/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87` should return an HTTP 200 status code with the JSON response:

```json
{
    "id": "cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87",
    "survey": "BRES"
}
```

### List Classifier Type Selectors

* Running the command `curl http://localhost:8080/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors` should return an HTTP 200 status code with the JSON response:

```json
[
  {
    "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
    "name": "COLLECTION_INSTRUMENT"
  },
  {
    "id": "e119ffd6-6fc1-426c-ae81-67a96f9a71ba",
    "name": "COMMUNICATION_TEMPLATE"
  }
]
```

### Get Classifier Types Selector

* Running the command `curl http://localhost:8080/surveys/cb0711c3-0ac8-41d3-ae0e-567e5ea1ef87/classifiertypeselectors/efa868fb-fb80-44c7-9f33-d6800a17c4da` should return an HTTP 200 status code with the JSON response:

```json
{
  "id": "efa868fb-fb80-44c7-9f33-d6800a17c4da",
  "name": "COLLECTION_INSTRUMENT",
  "classifierTypes": ["COLLECTION_EXERCISE", "RU_REF"]
}
```

## Testing
To follow once I've worked out how to write unit tests in Go :-)

Run the tests using:

```
make test
```

## Cleaning
To clobber the `build` directory tree, run:

```
make clean
```

## Copyright
Copyright (C) 2017 Crown Copyright (Office for National Statistics)
