FROM golang:alpine AS builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a -ldflags="-w -s" -tags netgo -o /go/bin/alerts .

FROM scratch
COPY --from=builder /go/bin/alerts /go/binalerts/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /go/bin
ADD data/db/migrations data/db/migrations

ENTRYPOINT ["./alerts"]
