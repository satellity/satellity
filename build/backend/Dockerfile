# syntax = docker/dockerfile:1-experimental

# Download deps
FROM --platform=${BUILDPLATFORM} golang:1.15.2-alpine AS base
WORKDIR /src
COPY go.* .
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Build executable bin
FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/app

# Support bin of different OS
FROM scratch AS bin-unix
COPY --from=build /out/app /

FROM bin-unix AS bin-darwin
FROM bin-unix AS bin-linux

FROM scratch AS bin-windows
COPY --from=build /out/app /app.exe

FROM bin-${TARGETOS} as bin

# Build docker image
FROM scratch AS dev
COPY --from=build /out/app /
COPY /internal/configs/config.example /internal/configs/config.yaml
EXPOSE 4000
CMD [ "/app" ]