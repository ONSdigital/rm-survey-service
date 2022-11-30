FROM golang:1.19.3-alpine3.16

EXPOSE 8080

RUN mkdir "/src"
WORKDIR "/src"

COPY main .

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

RUN go build
RUN ls

CMD "./main"