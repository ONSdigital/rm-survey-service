FROM ubuntu:18.04

RUN apt-get update\
     && apt-get install curl -y --no-install-recommends\
     && apt-get clean \
     && rm -rf /var/lib/apt/lists/*
EXPOSE 8080

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

ENTRYPOINT [ "/usr/local/bin/main" ]
