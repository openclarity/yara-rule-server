# syntax=docker/dockerfile:1@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021

FROM --platform=$BUILDPLATFORM golang:1.24rc2-bullseye@sha256:47104b6bfd9d29339a11fa498d82e73544f7dbe0c6f81dcccee2ac34bda8a080 AS builder

ARG TARGETPLATFORM

RUN --mount=type=cache,id=${TARGETPLATFORM}-apt,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
      gcc \
      libc6-dev

WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=.,target=/build,ro \
    go mod download -x

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
    go build -ldflags="-s -w -extldflags -static \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.Version=${VERSION}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.CommitHash=${COMMIT_HASH}' \
        -X 'github.com/openclarity/yara-rule-server/pkg/version.BuildTimestamp=${BUILD_TIMESTAMP}'" \
    -o /bin/yara-rule-server main.go

FROM alpine:3.20@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

RUN apk add --update --no-cache \
    util-linux \
    ca-certificates \
    libc6-compat

RUN apk --no-cache add  --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community yara=4.5.0-r0

COPY --from=builder --chmod=755 /bin/yara-rule-server /bin/yara-rule-server
COPY --link --chmod=644 config.example.yaml /etc/yara-rule-server/config.yaml

ENTRYPOINT ["/bin/yara-rule-server"]
