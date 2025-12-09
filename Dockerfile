FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制依赖文件并下载 (利用 Docker 缓存层)
COPY go.mod go.sum ./

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest

# 安装基础证书 (防止请求 HTTPS 报错)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]