[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c5adaae19b8f4b899ce935fe856a85d9)](https://www.codacy.com/app/sdcplatform/rm-survey-service?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ONSdigital/rm-survey-service&amp;utm_campaign=Badge_Grade) [![Docker Pulls](https://img.shields.io/docker/pulls/sdcplatform/surveysvc.svg)]()
[![Build Status](https://travis-ci.org/ONSdigital/rm-survey-service.svg?branch=master)](https://travis-ci.org/ONSdigital/rm-survey-service)

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
make docker
```

## Running
First compile the code using `make` then execute the binary in the background using `./surveysvc &` from within the `bin` directory within the `build` directory tree.

The following environment variables may be overridden:

| Environment Variable   | Purpose                                      | Default Value                                                   |
| :--------------------- | :------------------------------------------- | :-------------------------------------------------------------- |
| DATABASE_URL           | PostgreSQL *postgres* user connection string | postgres://postgres:password@localhost/postgres?sslmode=disable |
| PORT                   | HTTP listener port                           | :8080                                                           |
| security_user_name     | HTTP basic authentication user name          | N/A                                                             |
| security_user_password | HTTP basic authentication password           | N/A                                                             |
| CONN_MAX_LIFETIME      | Max lifetime of connection in pool in seconds| 0, so there is no time limit                                    |
| MAX_IDLE_CONN          | Max idle connections to have in pool         | 2                                                               |

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
Run the unit tests using:

```
make test
```

To run the integration tests using `make` run:
```bash
make integration-test
```

This will build a docker image from source and run the container then run all tests including integration tests.

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
