.PHONY: run_with_live_reload run migrate_up migrate_down migrate_status migrate_create migrate_reset build build_web

# 从 config.yaml 读取配置的辅助函数
# 如果 config.yaml 不存在，使用默认值
define get_config
$(shell grep "^$(1):" config.yaml 2>/dev/null | awk '{print $$2}' | tr -d '"' || echo "$(2)")
endef

# 动态读取数据库配置
DB_HOST := $(call get_config,DB_HOST,localhost)
DB_PORT := $(call get_config,DB_PORT,3306)
DB_USERNAME := $(call get_config,DB_USERNAME,root)
DB_PASSWORD := $(call get_config,DB_PASSWORD,root)
DB_DATABASE := $(call get_config,DB_DATABASE,go_api_template)

# 数据库连接字符串
DB_DSN := "$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)?parseTime=true"

# 开发运行
run_with_live_reload:
	air

run:
	go run cmd/main.go server

# 数据库迁移命令
migrate_up:
	@echo "Running migrations with DSN: $(DB_DSN)"
	goose -dir db/migrations mysql $(DB_DSN) up

migrate_down:
	@echo "Rolling back migration with DSN: $(DB_DSN)"
	goose -dir db/migrations mysql $(DB_DSN) down

migrate_status:
	@echo "Checking migration status with DSN: $(DB_DSN)"
	goose -dir db/migrations mysql $(DB_DSN) status

migrate_create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate_create NAME=your_migration_name"; \
		exit 1; \
	fi
	goose -dir db/migrations create $(NAME) go

migrate_reset:
	@echo "Resetting all migrations with DSN: $(DB_DSN)"
	goose -dir db/migrations mysql $(DB_DSN) reset

# 构建命令
build_web:
	@echo "Building web frontend..."
	@if [ -d "web" ]; then \
		cd web && npm run build && \
		echo "Copying build to internal/web..." && \
		rm -rf ../internal/web/build && \
		cp -r build ../internal/web/; \
	else \
		echo "⚠️  web 目录不存在,跳过前端构建"; \
		echo "如需嵌入前端,请创建 web 目录并放置前端项目"; \
	fi

build: build_web
	@echo "Building Go application..."
	go build -o go-api-template cmd/main.go
	@echo "✅ Build complete! Binary: ./go-api-template"
