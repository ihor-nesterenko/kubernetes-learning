# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

ENV GOPROXY=direct GOSUMDB=off
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest AS release
WORKDIR /
COPY --from=builder /app/app /app

# Add non root user for kubernetes container
RUN adduser -D -u 1001 -g "app-user" app-user

EXPOSE :8080
CMD ["/app"]