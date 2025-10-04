FROM golang:1.25-alpine

EXPOSE 8080

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

COPY build/darwin-arm64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

COPY ./cacert.pem /usr/local/share/ca-certificates/cacert.pem

RUN update-ca-certificates

RUN go build
RUN ls

CMD "./rm-survey-service"
