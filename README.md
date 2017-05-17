# Survey Service
This repository contains the Survey service. This microservice is a RESTful web service implemented using [Go](https://golang.org/). [API documentation](https://github.com/ONSdigital/rm-survey-service/blob/master/API.md).

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

Run `docker ps` and note the ID of the running Docker container. Initial data can be loaded into the PostgreSQL database by starting the Docker image and connecting, then loading the SQL files in the `sql` directory:

```
docker exec -it <container-id> /bin/sh
psql postgres://postgres:password@localhost/postgres?sslmode=disable
```

Manually copy and paste the contents of the SQL files in the following order:

```
groundzero.sql
survey_foundation_schema.sql
seed_data.sql
```

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
