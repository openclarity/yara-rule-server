# syntax=docker/dockerfile:1

FROM golang:1.21.5-alpine AS builder

RUN apk add --update --no-cache gcc g++ git ca-certificates build-base

ARG VERSION
ARG BUILD_TIMESTAMP
ARG COMMIT_HASH
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=.,target=/build,ro \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build \
    -ldflags="-s -w \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.Version=${VERSION}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.CommitHash=${COMMIT_HASH}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.BuildTimestamp=${BUILD_TIMESTAMP}'" \
    -o /bin/yara-rule-server main.go

FROM alpine:3.19

RUN apk --no-cache add util-linux yara ca-certificates

COPY --from=builder --chmod=755 /bin/yara-rule-server /bin/yara-rule-server
COPY --link --chmod=644 config.example.yaml /etc/yara-rule-server/config.yaml

ENTRYPOINT ["/bin/yara-rule-server"]
