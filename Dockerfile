FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build ./cmd/main.go

FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive
USER root

WORKDIR /rest
COPY --from=builder /build/main .

COPY . .

ENV APP_PORT "8081"

ENV REDIS_ADDR "redis:6379"
ENV REDIS_PASSWORD ""
ENV REDIS_DB 0
ENV REDIS_TIMER 15

ENV PSX_USER "admin"
ENV PSX_PASSWORD "admin"
ENV PSX_DBNAME "announcement"
ENV PSX_HOST "postgres"
ENV PSX_PORT 5432
ENV PSX_SSLMODE "disable"
ENV PSX_MAXCONNS 10
ENV PSX_TIMER 10

CMD ./main



