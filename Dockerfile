FROM golang:1.25-alpine

EXPOSE 8080

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

RUN go build
RUN ls

CMD "./rm-survey-service"