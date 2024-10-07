# 建構階段
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go clean -modcache

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 設定環境變數
ENV GO111MODULE=on

# 建構應用程式
RUN go build -o main .

# 運行階段
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

EXPOSE 3000

CMD ["./main"]
