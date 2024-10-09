# 建構階段
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth .

# 運行階段
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth .

EXPOSE 4000

CMD ["./auth"]
