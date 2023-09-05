# syntax=docker/dockerfile:1.2
FROM golang:1.20.7-alpine AS builder

RUN apk add --update --no-cache gcc g++ git ca-certificates build-base

# Copy code
COPY . /build
WORKDIR /build

ARG VERSION
ARG BUILD_TIMESTAMP
ARG COMMIT_HASH
ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build \
    -ldflags="-s -w \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.Version=${VERSION}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.CommitHash=${COMMIT_HASH}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.BuildTimestamp=${BUILD_TIMESTAMP}'" \
    -o bin/yara-rule-server main.go



FROM alpine:3.17

RUN apk upgrade
RUN apk add util-linux
RUN apk add --update yara --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing
RUN mkdir /etc/yara-rule-server

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/yara-rule-server ./yara-rule-server
COPY yara-rule-server-example /etc/yara-rule-server/config.yaml

WORKDIR /opt

ENTRYPOINT ["/app/yara-rule-server"]