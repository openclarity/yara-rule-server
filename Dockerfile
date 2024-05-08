# syntax=docker/dockerfile:1@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021

FROM --platform=$BUILDPLATFORM golang:1.22.3-bullseye@sha256:e72f9a1d29fbd6e1603df5a780e8f407caebef4dbb9f07536fc72f1c368298aa AS builder

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

FROM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

RUN apk add --update --no-cache \
    util-linux \
    ca-certificates \
    libc6-compat

RUN apk --no-cache add  --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community yara=4.5.0-r0

COPY --from=builder --chmod=755 /bin/yara-rule-server /bin/yara-rule-server
COPY --link --chmod=644 config.example.yaml /etc/yara-rule-server/config.yaml

ENTRYPOINT ["/bin/yara-rule-server"]
