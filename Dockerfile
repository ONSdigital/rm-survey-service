FROM alpine:latest
MAINTAINER John Topley "john.topley@ons.gov.uk"

# RUN apk update && apk upgrade && \
#     apk add --no-cache git openssh

COPY ./build/linux-amd64/bin/surveysvc /usr/local/bin/surveysvc

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/surveysvc"]
