FROM golang:1.19.3-alpine3.16

EXPOSE 8080

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

ENTRYPOINT [ "/usr/local/bin/main" ]

RUN go build
RUN ls