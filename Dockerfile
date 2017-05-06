FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Make the source code path
RUN mkdir -p /go/src/github.com/onsdigital/rm-survey-service/survey-api/

# Add all source code
ADD ./survey-api /go/src/github.com/onsdigital/rm-survey-service/survey-api/

# Run the Go installer
RUN go get -v -d github.com/onsdigital/rm-survey-service/survey-api/
RUN go install github.com/onsdigital/rm-survey-service/survey-api/

# Expose your port
EXPOSE 8080

# Indicate the binary as our entrypoint
ENTRYPOINT /go/bin/survey-api
