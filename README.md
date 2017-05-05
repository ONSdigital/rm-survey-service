# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/) and has the following responsibilities:

* Providing the classifier types applicable to a survey

## Prerequisites
* Install the [Go PostgreSQL driver]() using `go get github.com/lib/pq`
* Install the [Gorilla Mux URL router](http://www.gorillatoolkit.org/pkg/mux) using `go get github.com/gorilla/mux`

## Running
From `$GOPATH`, use `go run src/github.com/onsdigital/rm-survey-service/survey-api/main.go &` to start the Survey service in the background. The following environment variables may be overridden:

<table>
  <thead>
    <tr>
      <th>Environment Variable</th>
      <th>Purpose</th>
      <th>Default Value</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>SURVEY_SERVICE_DB_CONNECT</td>
      <td>PostgreSQL database connection string</td>
      <td>postgres://postgres:password@localhost/postgres?sslmode=disable</td>
    </tr>
    <tr>
      <td>SURVEY_SERVICE_PORT</td>
      <td>HTTP listener port</td>
      <td>:8080</td>
  </tbody>
</table>

## API Examples

* Running the command `curl http://localhost:8080/surveys` should return an HTTP 200 status code with the JSON response:

```json
["BRES"]
```

* Running the command `curl http://localhost:8080/surveys/bres/classifiertypes` should return an HTTP 200 status code with the JSON response:

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

## Copyright
Copyright (C) 2017 Crown Copyright (Office for National Statistics)
