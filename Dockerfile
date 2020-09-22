FROM golang:1.15 AS builder

WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o serverd ./cmd/...

# FROM alpine:3.7
# COPY --from=builder /app .
