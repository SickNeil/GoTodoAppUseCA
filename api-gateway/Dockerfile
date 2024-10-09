# 建構階段：使用官方的 Go 圖像作為基礎映像
FROM golang:1.23-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製源代碼到容器中
COPY . .

# 構建 Go 應用程式
RUN go build -o api-gateway .

# 運行階段：使用輕量級的 Alpine 映像
FROM alpine:latest

# 設置工作目錄
WORKDIR /root/

# 複製建構階段的構建成果到運行階段
COPY --from=builder /app/api-gateway .

# 暴露服務端口
EXPOSE 5000

# 設置容器啟動時運行的命令
CMD ["./api-gateway"]
