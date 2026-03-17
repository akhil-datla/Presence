# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /presence .

# Stage 2: Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S presence \
    && adduser -S presence -G presence

RUN mkdir -p /data && chown presence:presence /data
VOLUME /data

ENV DATABASE_PATH=/data/presence.db

COPY --from=builder /presence /presence

EXPOSE 8080

USER presence

ENTRYPOINT ["/presence"]
