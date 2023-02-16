# syntax=docker/dockerfile:1

# We use a multi-stage build setup.
# (https://docs.docker.com/build/building/multi-stage/)

# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.19.5 AS builder
# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

COPY . /go/src/github.com/ghowland/sireus/
WORKDIR /go/src/github.com/ghowland/sireus/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o build/sireus code/sireus.go

#RUN go test -cover -v ./...

# Stage 2 (to create a downsized "container executable", ~5MB)

# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:
#
#   FROM alpine:3.17.1
#   RUN apk --no-cache add ca-certificates
#
#FROM scratch
FROM ubuntu:22.04
WORKDIR /root/
COPY --from=builder /go/src/github.com/ghowland/sireus/build/sireus .

RUN mkdir ./config
RUN mkdir ./web
RUN mkdir ./static_web

# Copy configuration files needed
COPY --from=builder /go/src/github.com/ghowland/sireus/config/ ./config/
COPY --from=builder /go/src/github.com/ghowland/sireus/web/ ./web/
COPY --from=builder /go/src/github.com/ghowland/sireus/static_web/ ./static_web/

# Install Prometheus, for the demo
RUN /usr/bin/apt update
RUN /usr/bin/apt install -y prometheus
RUN /usr/bin/apt install -y net-tools

COPY --from=builder /go/src/github.com/ghowland/sireus/demo/config/prometheus.yml /etc/prometheus/prometheus.yml

EXPOSE 3000
EXPOSE 8611
EXPOSE 8709
EXPOSE 9090
CMD service prometheus start && ./sireus
