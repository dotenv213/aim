FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN go build -mod=vendor -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081 50051

CMD ["./main"]