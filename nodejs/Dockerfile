# syntax=docker/dockerfile:1

FROM node:22-bookworm-slim AS deps

WORKDIR /app

RUN  \
  --mount=type=bind,source=package.json,target=package.json \
  --mount=type=bind,source=package-lock.json,target=package-lock.json \
  npm ci --omit=dev

FROM node:22-bookworm-slim AS build

WORKDIR /app

RUN --mount=type=cache,target=node_modules,sharing=locked \
  --mount=type=bind,source=package.json,target=package.json \
  --mount=type=bind,source=package-lock.json,target=package-lock.json \
  npm ci

RUN --mount=type=cache,target=node_modules,sharing=locked \
  --mount=type=bind,source=package.json,target=package.json \
  --mount=type=bind,source=main.ts,target=main.ts \
  npm run build

FROM node:22-bookworm-slim AS final

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
  --mount=type=cache,target=/var/cache/apt,sharing=locked \
  apt-get update && apt-get install -y ca-certificates tini

USER nobody

ENV NODE_ENV production

COPY --from=deps /app/node_modules /app/node_modules
COPY --from=build /app/main.js /app/main.js

ENTRYPOINT [ "/usr/bin/tini", "--" ]
CMD [ "node", "/app/main.js" ]
