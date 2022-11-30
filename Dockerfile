FROM golang:1.19.3-alpine3.16

EXPOSE 8080

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN go build
RUN ls