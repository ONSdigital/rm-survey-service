FROM ubuntu:18.04

RUN apt-get update && apt-get install curl -y
EXPOSE 8080

COPY build/linux-amd64/bin/main /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/main" ]
