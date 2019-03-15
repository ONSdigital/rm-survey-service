FROM ubuntu:18.04

RUN apt-get update\
     && apt-get clean \
     && rm -rf /var/lib/apt/lists/*
EXPOSE 8080

RUN groupadd --gid 995 surveysvc && \
    useradd --create-home --system --uid 995 --gid surveysvc surveysvc
USER surveysvc

COPY build/linux-amd64/bin/main /usr/local/bin/

COPY db-migrations /db-migrations

ENTRYPOINT [ "/usr/local/bin/main" ]
