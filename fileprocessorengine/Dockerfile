FROM golang:1.18-alpine AS build-env

COPY ./go.* /src/
COPY ./fileprocessorengine /src/fileprocessorengine
COPY ./model /src/model
COPY ./pkg /src/pkg
COPY ./utils /src/utils

WORKDIR /src/fileprocessorengine
RUN go build -o /fileprocessorengine


ENTRYPOINT ["/fileprocessorengine"]
