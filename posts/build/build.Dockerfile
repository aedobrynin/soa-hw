FROM golang:1.21.1-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

WORKDIR /app/cmd/posts
RUN go build -o main