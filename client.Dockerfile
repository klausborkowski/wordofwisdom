FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o client ./cmd/client

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/client .

ENTRYPOINT ["./client"]
