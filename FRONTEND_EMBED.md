# 前端资源嵌入说明

本模板支持将前端静态资源嵌入到 Go 二进制文件中,实现单一可执行文件部署。

## 快速开始

### 1. 创建前端项目

在项目根目录创建 `web` 目录,并放置你的前端项目(React/Vue/任何前端框架):

```bash
# 示例: 创建 React 项目
npx create-react-app web
cd web
npm install react-router-dom
```

### 2. 配置前端路由

如果使用 React Router,参考 ai-gateway 项目的实现:
- 创建 Header 组件(包含导航菜单)
- 创建页面组件(Dashboard, Users 等)
- 在 App.js 中配置路由

### 3. 构建和部署

```bash
# 完整构建(前端 + 后端)
make build

# 运行
./go-api-template server -c config.yaml
```

访问 http://localhost:3000 即可看到前端应用。

## 开发模式

开发时前后端分离运行:

```bash
# 终端 1: 启动后端
make run

# 终端 2: 启动前端开发服务器
cd web && npm start
```

前端 `package.json` 中配置 proxy 代理 API 请求到后端:

```json
{
  "proxy": "http://localhost:3000"
}
```

## 目录结构

```
.
├── cmd/
│   ├── main.go              # 引用 web.BuildFS 和 web.IndexPage
│   └── server/
│       ├── cmd.go           # 配置路由,调用 SetWebRouter
│       └── web_router.go    # Web 路由配置
├── internal/
│   ├── common/
│   │   └── embed_fs.go      # embed 文件系统工具
│   └── web/
│       ├── embed.go         # embed 声明
│       ├── .gitkeep         # git 占位文件
│       └── build/           # 前端构建产物(gitignore)
├── web/                     # 前端项目(可选)
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── build/               # 前端构建输出
├── Makefile                 # 包含 build 和 build_web 命令
└── .gitignore              # 忽略构建产物
```

## 路由约定

为避免冲突,遵循以下路由约定:

- **`/`** - 前端首页
- **`/users`, `/settings` 等** - 前端路由
- **`/api/*`** - 所有后端 API
- **`/api/health`** - 健康检查(原 `/` 路由移至此)

## 技术细节

### 1. Embed 声明 (internal/web/embed.go)

```go
//go:embed build
var BuildFS embed.FS

//go:embed build/index.html
var IndexPage []byte
```

### 2. SPA 路由处理 (cmd/server/web_router.go)

NoRoute 中间件处理:
- `/api/*` - 返回 404 API 错误
- `/static/*`, `/assets/*` - 静态资源,走文件系统
- 其他路径 - 返回 index.html,交给前端路由

### 3. 构建流程

`make build` 执行:
1. `make build_web` - 构建前端(如果 web 目录存在)
2. 复制 `web/build/` 到 `internal/web/build/`
3. Go 编译时嵌入 `internal/web/build/`

### 4. 无前端项目时

如果不创建 `web` 目录:
- `make build` 会跳过前端构建
- `internal/web/build/` 为空
- embed 会报错,需要创建空的占位文件

解决方案: `internal/web/.gitkeep` 确保目录存在。

## 示例项目

参考 `ai-gateway` 项目查看完整实现:
- `/Users/mouyong/Code/mouyong/ai-gateway/web/` - React 前端
- 导航栏、路由、API 集成示例

## 常见问题

### 1. 编译时找不到 build 目录

确保先运行 `make build_web` 或创建空的 `internal/web/build/index.html`。

### 2. 修改前端不生效

- 开发模式: 使用 `cd web && npm start`,访问 3000 端口
- 生产模式: 重新 `make build`

### 3. API 请求 404

检查:
- 后端路由是否以 `/api` 开头
- 前端 package.json 中 proxy 配置是否正确
- NoRoute 配置是否正确处理 `/api/*`

## 性能优化

### 前端打包优化

- 使用代码分割(React.lazy)
- 启用 tree shaking
- 分析打包体积: `npm run build -- --stats`

### Gzip 压缩

在 `web_router.go` 中添加:

```go
import "github.com/gin-contrib/gzip"

router.Use(gzip.Gzip(gzip.DefaultCompression))
```

## 部署

### Docker 部署

Dockerfile 示例:

```dockerfile
# 构建阶段
FROM node:18 AS web-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

FROM golang:1.23 AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /app/web/build ./internal/web/build
RUN go build -o go-api-template cmd/main.go

# 运行阶段
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=go-builder /app/go-api-template .
COPY --from=go-builder /app/config.example.yaml ./config.yaml
CMD ["./go-api-template", "server"]
```

### 单文件部署

```bash
make build
scp go-api-template config.yaml server:/path/to/app/
ssh server 'cd /path/to/app && ./go-api-template server'
```
