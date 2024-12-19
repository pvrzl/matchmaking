FROM golang:1.21-alpine

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOSE=linux

WORKDIR /app

COPY . .

# Downloading dependencies
RUN go mod download && go mod verify 

# Install build tools
RUN apk update && \
    apk add --no-cache make

# Install CompileDaemon
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="make build-http" -command="./app-http"
