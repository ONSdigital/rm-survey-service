FROM scratch
MAINTAINER John Topley "john.topley@ons.gov.uk"
EXPOSE 8080
COPY surveysvc /
ENTRYPOINT ["/surveysvc"]
