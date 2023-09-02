FROM golang:1.18-alpine AS build-env

COPY ./go.* /src/
COPY ./clientsvc /src/clientsvc
COPY ./model /src/model
COPY ./pkg /src/pkg
COPY ./utils /src/utils

WORKDIR /src/clientsvc
RUN go build -o /clientsvc


ENTRYPOINT ["/clientsvc"]
