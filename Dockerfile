FROM --platform=$BUILDPLATFORM golang:1.25-alpine

EXPOSE 8080

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then \
      cp build/linux-amd64/bin/main /usr/local/bin/main; \
    elif [ "$BUILDPLATFORM" = "darwin/arm64" ]; then \
      cp build/darwin-arm64/bin/main /usr/local/bin/main; \
    fi

COPY db-migrations /db-migrations

RUN go build
RUN ls

CMD "./rm-survey-service"
