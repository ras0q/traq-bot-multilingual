# syntax=docker/dockerfile:1

FROM denoland/deno:debian-1.46.3 AS base

WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/deno,sharing=locked \
  --mount=type=bind,target=. \
  deno cache .

FROM base AS build

RUN --mount=type=cache,target=/root/.cache/deno,sharing=locked \
  --mount=type=bind,target=. \
  deno compile --allow-env --allow-net --unstable-net --output /main main.ts

FROM debian:bookworm-slim AS final

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
  --mount=type=cache,target=/var/cache/apt,sharing=locked \
  apt-get update && apt-get install -y ca-certificates

USER nobody

COPY --from=build /main /main

ENTRYPOINT [ "/main" ]
