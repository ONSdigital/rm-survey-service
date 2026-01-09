FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN go build -v -o main
RUN chmod 755 main

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS final-stage

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
