FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3333

ENTRYPOINT ["./server"]
