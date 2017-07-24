[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c5adaae19b8f4b899ce935fe856a85d9)](https://www.codacy.com/app/sdcplatform/rm-survey-service?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ONSdigital/rm-survey-service&amp;utm_campaign=Badge_Grade) [![Docker Pulls](https://img.shields.io/docker/pulls/sdcplatform/surveysvc.svg)]()

# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/). This service features structured JSON logging, a self-bootstrapping database schema and database connection code that retries the connection if it's not available, increasing the time between each attempt. This [eliminates the need to deploy services in a specific order](https://medium.com/@kelseyhightower/12-fractured-apps-1080c73d481c).

## Prerequisites
* Install the [Godep](https://github.com/tools/godep) package manager using `go get github.com/tools/godep`
* Run `godep get` to download and install the other dependencies managed by Godep

## API
See [API.md](https://github.com/ONSdigital/rm-survey-service/blob/master/API.md) for API documentation.

## Building
Install Go and ensure your `GOPATH` environment variable is set (usually it's `~/go`).

### Make
A Makefile is provided for compiling the code:

```
make
```

The compiled executable is placed within the `build` directory tree.

### Docker Image
To build the Docker image, first compile the code using `make` then from the project root run:

```
docker build -t "sdcplatform/surveysvc" .
```

## Running
First compile the code using `make` then execute the binary in the background using `./surveysvc &` from within the `bin` directory within the `build` directory tree.

The following environment variables may be overridden:

| Environment Variable | Purpose                                      | Default Value                                                   |
| :------------------- | :------------------------------------------- | :-------------------------------------------------------------- |
| DATABASE_URL         | PostgreSQL *postgres* user connection string | postgres://postgres:password@localhost/postgres?sslmode=disable |
| PORT                 | HTTP listener port                           | :8080                                                           |

### Docker Image and PostgreSQL
To start Docker containers for both PostgreSQL and the Survey service, run:

```
docker-compose up -d
```

To stop and remove the two Docker containers, run:

```
docker-compose down
```

## Testing
To follow once I've worked out how to write unit tests in Go :-)

Run the tests using:

```
make test
```

## Deployment
To deploy to Cloud Foundry, run one of the targets below depending on the Cloud Foundry space you wish to push to:

```
make push-ci
make push-demo
make push-dev
make push-int
make push-test
```

## Cleaning
To clobber the `build` directory tree that's created when running `make`, run:

```
make clean
```

## Copyright
Copyright (C) 2017 Crown Copyright (Office for National Statistics)
