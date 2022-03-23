# Container for building Go binary.
FROM golang:1.17-alpine AS builder
# Cgo support
RUN apk add --no-cache build-base
# Copy app.
WORKDIR /app
COPY . .
# Build with Go module and Go build caches.
RUN \
   --mount=type=cache,target=/go/pkg \
   --mount=type=cache,target=/root/.cache/go-build \
   go build -o charon main.go

# Copy final binary into light stage.
FROM alpine:3
ARG GITHUB_SHA=local
ENV GITHUB_SHA=${GITHUB_SHA}
COPY --from=builder /app/charon /usr/local/bin/
# Don't run container as root
ENV USER=charon
ENV UID=12345
ENV GID=23456
RUN addgroup -g "$GID" "$USER"
RUN adduser \
    --disabled-password \
    --gecos "charon" \
    --home "/opt/$USER" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
RUN chown charon /usr/local/bin/charon
RUN chmod u+x /usr/local/bin/charon
WORKDIR "/opt/$USER"
USER charon
CMD ["/usr/local/bin/charon","run"]
# Used by GitHub to associate container with repo.
LABEL org.opencontainers.image.source="https://github.com/obolnetwork/charon"
