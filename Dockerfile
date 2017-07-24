FROM alpine:3.6
MAINTAINER John Topley "john.topley@ons.gov.uk"
RUN apk update && apk upgrade
COPY build/linux-amd64/bin/surveysvc /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/surveysvc"]
