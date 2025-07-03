# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/monitor && go build -o /monitor

FROM alpine:3.19
COPY --from=builder /monitor /monitor
EXPOSE 8080
ENTRYPOINT ["/monitor"] 