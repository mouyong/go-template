# Go Template

基于 Gin + GORM 的 Go API 项目模板，支持可选的数据库和 RabbitMQ 配置。

## 快速开始

### 1. 配置文件

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml` 配置数据库和 RabbitMQ（可选）。

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 安装数据库迁移工具

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. 配置数据库连接

编辑 `config.yaml` 配置数据库：

```yaml
DB_HOST: "localhost"
DB_PORT: 3306
DB_DATABASE: "go_template"
DB_USERNAME: "root"
DB_PASSWORD: "root"
```

**注意**: Makefile 会自动从 `config.yaml` 读取数据库配置，无需单独配置。

### 5. 运行数据库迁移

```bash
# 查看迁移状态
make migrate_status

# 执行迁移
make migrate_up

# 回滚上一次迁移
make migrate_down

# 重置所有迁移
make migrate_reset
```

### 6. 启动服务

```bash
# 使用 air 热更新（推荐开发环境）
make run_with_live_reload

# 或直接运行
make run
```

服务将在 `http://localhost:3000` 启动。

## API 路由

- `GET /` - 健康检查
- `GET /api/hello?name=World` - Hello 示例
- `POST /api/echo` - Echo 示例（JSON 回显）

## 数据库迁移

### 创建新迁移

```bash
make migrate_create NAME=create_products_table
```

这将在 `db/migrations` 目录创建一个新的迁移文件。

### 迁移文件格式

使用 Go 文件迁移（支持跨数据库）：

```go
package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upCreateProductsTable, downCreateProductsTable)
}

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

func upCreateProductsTable(ctx context.Context, tx *sql.Tx) error {
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&Product{})
}

func downCreateProductsTable(ctx context.Context, tx *sql.Tx) error {
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		return err
	}
	return db.Migrator().DropTable(&Product{})
}
```

**优势**: 使用 GORM 的 AutoMigrate，支持 MySQL、PostgreSQL、SQLite 等多种数据库。

## 目录结构

```
.
├── cmd/                    # 命令行、入口目录
│   ├── main.go             # 入口文件
│   └── server/             # server 命令
│       └── cmd.go          # 服务启动、路由配置
├── config.yaml             # 配置文件
├── config.example.yaml     # 配置文件示例
├── db/                     # 数据库相关
│   └── migrations/         # 数据库迁移文件
├── internal/               # 内部代码
│   ├── handlers/           # HTTP 处理器
│   │   ├── common.go       # 公共响应结构
│   │   └── example.go      # 示例处理器
│   ├── initialization/     # 初始化
│   │   ├── config.go       # 配置加载
│   │   └── db.go           # 数据库初始化
│   └── models/             # 数据模型
├── pkg/                    # 可复用的包
│   └── rabbitmq/           # RabbitMQ 客户端
│       └── client.go
├── Makefile                # Make 命令
└── README.md
```

## 配置说明

### 服务配置

```yaml
HTTP_HOST: 0.0.0.0
HTTP_PORT: 3000
```

### 数据库配置（可选）

如果不配置数据库，服务仍可正常启动。

```yaml
DB_HOST: "localhost"
DB_PORT: 3306
DB_DATABASE: "go_api_template"
DB_USERNAME: "root"
DB_PASSWORD: "root"
```

### RabbitMQ 配置（可选）

如果不配置 RabbitMQ，服务仍可正常启动。

```yaml
MQ_HOST: "localhost"
MQ_PORT: 5672
```

## 开发

### 热更新开发

使用 air 进行热更新开发：

```bash
air
```

配置文件: `.air.toml`

### 添加新路由

在 `cmd/server/cmd.go` 中添加：

```go
r.GET("/api/your-route", handlers.YourHandler)
```

### 添加新 Handler

在 `internal/handlers/` 中创建新文件：

```go
func YourHandler(c *gin.Context) {
    resp := NewResp(c)
    resp.successWithData(gin.H{
        "message": "success",
    }, nil)
}
```

### RabbitMQ 队列监听

在 `pkg/rabbitmq/client.go` 的 `ListenQueue()` 函数中：

```go
func ListenQueue() {
    if RabbitmqChannel == nil {
        return
    }

    // 启动自定义队列监听
    StartQueue("your_queue", YourHandler)
}

func YourHandler(body []byte) error {
    // 处理消息
    log.Printf("Processing: %s", string(body))
    return nil
}
```

## Make 命令

```bash
make run                    # 运行服务
make run_with_live_reload   # 热更新运行
make migrate_up             # 执行数据库迁移
make migrate_down           # 回滚迁移
make migrate_status         # 查看迁移状态
make migrate_create NAME=xxx # 创建新迁移
make migrate_reset          # 重置所有迁移
```

## License

MIT
