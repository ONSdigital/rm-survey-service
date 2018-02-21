FROM ubuntu:18.04

RUN apt update && apt install curl -y
EXPOSE 8080

COPY build/linux-amd64/bin/main /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/main" ]
