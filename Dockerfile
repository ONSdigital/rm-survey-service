FROM golang:alpine
MAINTAINER John Topley "john.topley@ons.gov.uk"

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

COPY ./models /go/src/github.com/onsdigital/rm-survey-service/models/
COPY ./surveysvc.go /go/src/github.com/onsdigital/rm-survey-service/

RUN go get -v -d github.com/onsdigital/rm-survey-service/
RUN go build -o "$GOPATH"/bin/surveysvc github.com/onsdigital/rm-survey-service/

EXPOSE 8080
ENTRYPOINT ["/go/bin/surveysvc"]
