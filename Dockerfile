# syntax=docker/dockerfile:1

ARG GO_VERSION="1.26"
ARG GOLANGCI_LINT_VERSION="v2.11.3"

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine AS golangci-lint

FROM golang:${GO_VERSION}-alpine AS base
ENV GOFLAGS="-buildvcs=false"
WORKDIR /src

FROM base AS test
RUN --mount=type=bind,target=.,rw \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build <<EOT
  set -ex
  go test -v -coverprofile=/tmp/coverage.txt -covermode=atomic ./...
  go tool cover -func=/tmp/coverage.txt
EOT

FROM scratch AS test-coverage
COPY --from=test /tmp/coverage.txt /

FROM base AS lint
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/root/.cache \
    --mount=from=golangci-lint,source=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
  golangci-lint run ./...

FROM base AS build
ARG VERSION="dev"
ENV CGO_ENABLED=0
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
  go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/gomod-updates ./cmd/gomod-updates

FROM scratch AS binary
COPY --from=build /out/gomod-updates /

FROM golang:${GO_VERSION}-alpine AS image
RUN apk add --no-cache git
COPY --from=build /out/gomod-updates /usr/bin/gomod-updates
ENTRYPOINT ["gomod-updates"]
