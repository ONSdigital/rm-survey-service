# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/). This service features structured JSON logging, a self-bootstrapping database schema and database connection code that retries the connection if it's not available, increasing the time between each attempt. This [eliminates the need to deploy services in a specific order](https://medium.com/@kelseyhightower/12-fractured-apps-1080c73d481c).

## Prerequisites
* Install the [Godep](https://github.com/tools/godep) package manager using `go get github.com/tools/godep`
* Run `godep get` to download and install the other dependencies managed by Godep

## API
See [API.md](https://github.com/ONSdigital/rm-survey-service/blob/main/API.md) for API documentation.

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

| Environment Variable   | Purpose                                       | Default Value                                                   |
|:-----------------------|:----------------------------------------------|:----------------------------------------------------------------|
| DATABASE_URL           | PostgreSQL *postgres* user connection string  | postgres://postgres:password@localhost/postgres?sslmode=disable |
| PORT                   | HTTP listener port                            | :8080                                                           |
| security_user_name     | HTTP basic authentication user name           | N/A                                                             |
| security_user_password | HTTP basic authentication password            | N/A                                                             |
| CONN_MAX_LIFETIME      | Max lifetime of connection in pool in seconds | 0, so there is no time limit                                    |
| MAX_IDLE_CONN          | Max idle connections to have in pool          | 2                                                               |


### Running locally with `go run`

It's possible to run the application locally using just `go run main.go` but in order to do so there are a few things
to change:
- Update the `dataSource` database url in `main.go` to a valid postgres instance
- Update the `migrationSource` variable in `main.go` to `file://./db-migrations`

Once the app is running you should call the app with blank http auth values. This is because it's looking for the
`security_user_name` and `security_user_password` environment variables.  If they're not set, they just default to nothing.

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


## Cleaning
To clobber the `build` directory tree that's created when running `make`, run:

```
make clean
```

## Database Migrations in Kubernetes
Due to limitations on the DB migration library in use in this project the following steps need to be followed to ensure
a database migration is successful on a Kubernetes environment.

1. deploy the application via Spinnaker
1. scale the replicas to 1 pod (so it does the migration)
1. scale the replicas to 0 pods (to release the lock)
1. scale the replicas back to 2 pods again (or redeploy via Spinnaker)

## Copyright
Copyright (C) 2017 - 2020 Crown Copyright (Office for National Statistics)
build please