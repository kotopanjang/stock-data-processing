FROM golang:1.18-alpine AS build-env

COPY ./go.* /src/
COPY ./dataprocessorengine /src/dataprocessorengine
COPY ./model /src/model
COPY ./pkg /src/pkg
COPY ./utils /src/utils

WORKDIR /src/dataprocessorengine
RUN go build -o /dataprocessorengine


ENTRYPOINT ["/dataprocessorengine"]
