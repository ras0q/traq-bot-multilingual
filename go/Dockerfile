# syntax=docker/dockerfile:1

FROM golang:1.23.1-bookworm AS base

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

FROM base AS build

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 go build -o /main .

FROM debian:bookworm-slim AS final

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
  --mount=type=cache,target=/var/cache/apt,sharing=locked \
  apt-get update && apt-get install -y ca-certificates

USER nobody

COPY --from=build /main /main

ENTRYPOINT [ "/main" ]
