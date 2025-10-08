# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server .

FROM alpine:3.19
WORKDIR /app

RUN adduser -D -g '' appuser \
    && apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /app/server
COPY --from=builder /app/generation/config /app/generation/config

EXPOSE 3000
ENV PORT=3000
USER appuser

ENTRYPOINT ["/app/server"]
