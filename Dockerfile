FROM ubuntu:latest

COPY build/linux-amd64/bin/surveysvc /usr/local/bin/

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/surveysvc" ]
