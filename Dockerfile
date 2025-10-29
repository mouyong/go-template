# build stage
FROM node:20-slim AS frontend-builder
WORKDIR /app
RUN yarn global add bun --registry=https://registry.npmmirror.com
COPY ./web/ .
RUN bun install --registry=https://registry.npmmirror.com
# RUN cp .env.production .env
RUN bun run build

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

# 复制前端构建结果
COPY --from=frontend-builder /app/build ./internal/web/build

# 验证前端文件
RUN ls -la ./internal/web/build/ && \
    test -f ./internal/web/build/index.html && \
    echo "✓ Frontend files copied successfully"

# 构建二进制文件
RUN go build -ldflags="-s -w" -o app cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制二进制文件
COPY --from=builder /app/app .

# 复制配置文件模板（可选）
COPY config.production.yaml ./config.yaml

EXPOSE 3000

CMD ["./app", "server"]
