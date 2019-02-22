FROM ubuntu:18.04

RUN apt-get update\
     && apt-get clean \
     && rm -rf /var/lib/apt/lists/*
EXPOSE 8080

RUN groupadd -g 995 surveysvc && \
    useradd -r -u 995 -g surveysvc surveysvc
USER surveysvc

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

ENTRYPOINT [ "/usr/local/bin/main" ]
