# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/). [API documentation](https://github.com/ONSdigital/rm-survey-service/blob/master/API.md).

## Prerequisites
* Install the [Go PostgreSQL driver](https://github.com/lib/pq) using `go get github.com/lib/pq`
* Install the [Echo web framework](https://github.com/labstack/echo) using `go get -u github.com/labstack/echo/...`

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

| Environment Variable | Purpose                                      | Default Value                                                   |
| :------------------- | :------------------------------------------- | :-------------------------------------------------------------- |
| PG_DATABASE_URL      | PostgreSQL *postgres* user connection string | postgres://postgres:password@localhost/postgres?sslmode=disable |
| DATABASE_URL         | PostgreSQL *survey* user connection string   | postgres://survey:password@localhost/postgres?sslmode=disable   |
| PORT                 | HTTP listener port                           | :8080                                                           |

### Docker Image and PostgreSQL
To start the Docker image, run:

```
docker-compose up -d
```

Run `docker ps` and note the ID of the running Docker container. Initial data can be loaded into the PostgreSQL database by starting the Docker image and connecting, then loading the `sql/bootstrap.sql` file:

```
docker exec -it <container-id> /bin/sh
psql postgres://postgres:password@localhost/postgres?sslmode=disable
```

Manually copy and paste the contents of the `sql/bootstrap.sql` file.

## API
See [API.md](https://github.com/ONSdigital/rm-survey-service/blob/master/API.md) for API documentation.

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
