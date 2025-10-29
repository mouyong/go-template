# # 前端构建阶段
# FROM oven/bun:1-alpine AS frontend-builder

# WORKDIR /app

# # 复制前端代码
# COPY web/package.json web/bun.lock ./
# RUN bun install

# COPY web/ ./
# RUN bun run build

# Go 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 设置国内代理加速
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# # 复制前端构建结果
# COPY --from=frontend-builder /app/build ./internal/web/build

# 构建二进制文件
RUN go build -ldflags="-s -w" -o app cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制二进制文件
COPY --from=builder /app/app .

# 复制配置文件模板（可选）
COPY config.production.yaml ./config.yaml

EXPOSE 3000

CMD ["./app", "server"]
