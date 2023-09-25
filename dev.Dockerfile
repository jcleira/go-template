FROM golang:alpine AS builder

# On this dev.Dockerfile I want to keep a build folder with the codebase, but I
# will also build and freeze a compiled version into /go/bin (that will include
# runtime dependencies. Ex: migrations.

# Building & codebase folder
RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a -ldflags="-w -s" -tags netgo -o /go/bin/alerts .

# Runtime folder configuration
WORKDIR /go/bin
ADD data/db/migrations data/db/migrations

# We are not defining an entrypoint or command here, that should be in the
# docker-compose configuration as this service migth be started with multiple
# commands, Ex: workers, jobs, web servers, etc.
