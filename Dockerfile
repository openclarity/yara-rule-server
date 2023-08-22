# syntax=docker/dockerfile:1.2
FROM golang:1.20.7-alpine AS builder

RUN apk add --update --no-cache gcc g++ git ca-certificates build-base

# Copy code
COPY . /build
WORKDIR /build

ARG BUILD_TIMESTAMP
ARG IMAGE_VERSION

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build \
    -ldflags="-s -w" \
    -o bin/yara-rule-server main.go

FROM ubuntu:22.04

RUN apt update
RUN apt install -y util-linux ca-certificates yara
ADD https://raw.githubusercontent.com/Yara-Rules/rules/master/index_gen.sh /usr/local/bin
RUN chmod 755 /usr/local/bin/index_gen.sh

WORKDIR /app

COPY --from=builder /build/bin/yara-rule-server ./yara-rule-server

WORKDIR /opt

ENTRYPOINT ["/app/yara-rule-server"]