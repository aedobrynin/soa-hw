FROM golang:1.21.1-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

WORKDIR /app/cmd/core
RUN go build -o main

# TODO: run oapi-codegen here
# TODO: run posts-client-generate here
