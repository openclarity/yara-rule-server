# syntax=docker/dockerfile:1.2
FROM golang:1.21.0-alpine AS builder

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

FROM alpine:3.18

RUN apk upgrade
RUN apk add util-linux

WORKDIR /app

COPY --from=builder /build/bin/yara-rule-server ./yara-rule-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/yara-rule-server"]