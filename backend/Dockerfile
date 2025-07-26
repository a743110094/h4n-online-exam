# 使用官方 Golang 镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包（包括CGO依赖）
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（启用CGO支持SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main .

# 使用轻量级的 alpine 镜像作为运行环境
FROM alpine:latest

# 安装 ca-certificates（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制环境变量文件（可选）
COPY --from=builder /app/.env.example .env

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]