FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage
ARG TARGETOS
ARG TARGETARCH

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -o main
RUN chmod 755 main

FROM alpine:3.24 AS final-stage

RUN apk add --no-cache ca-certificates

RUN addgroup -S survey-group && adduser -S survey-user -G survey-group
RUN mkdir -p "/opt/survey"
RUN chown survey-user:survey-group /opt/survey

WORKDIR "/opt/survey"
COPY --from=build-stage /src/main .
COPY --from=build-stage /src/db-migrations /db-migrations

RUN chmod 550 /opt/survey/main
RUN chown survey-user:survey-group /opt/survey/main
RUN chown survey-user:survey-group /db-migrations

USER survey-user

CMD "./main"
