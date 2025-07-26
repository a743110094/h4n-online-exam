# 在线考试系统后端 Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop deps fmt lint

# 默认目标
help:
	@echo "在线考试系统后端 - 可用命令:"
	@echo "  build        - 构建应用程序"
	@echo "  run          - 运行应用程序"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  deps         - 安装/更新依赖"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  docker-build - 构建 Docker 镜像"
	@echo "  docker-run   - 使用 Docker Compose 运行服务"
	@echo "  docker-stop  - 停止 Docker 服务"
	@echo "  docker-logs  - 查看 Docker 日志"
	@echo "  db-init      - 初始化数据库"

# 构建应用程序
build:
	@echo "构建应用程序..."
	go build -o bin/main .

# 运行应用程序
run:
	@echo "启动应用程序..."
	go run main.go

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 安装/更新依赖
deps:
	@echo "安装/更新依赖..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	go vet ./...

# 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t online-exam-backend .

# 使用 Docker Compose 运行服务
docker-run:
	@echo "启动 Docker 服务..."
	docker-compose up -d

# 停止 Docker 服务
docker-stop:
	@echo "停止 Docker 服务..."
	docker-compose down

# 查看 Docker 日志
docker-logs:
	@echo "查看应用日志..."
	docker-compose logs -f backend

# 初始化数据库
db-init:
	@echo "初始化数据库..."
	@if command -v psql >/dev/null 2>&1; then \
		psql -h localhost -U postgres -f scripts/init_db.sql; \
	else \
		echo "请安装 PostgreSQL 客户端或使用 Docker 运行数据库"; \
	fi

# 开发环境快速启动
dev: deps
	@echo "启动开发环境..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "已创建 .env 文件，请根据需要修改配置"; \
	fi
	make run

# 生产环境构建
prod-build:
	@echo "生产环境构建..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/main .

# 安装开发工具
tools:
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 生成 API 文档
docs:
	@echo "生成 API 文档..."
	swag init

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	go run main.go --migrate

# 创建新的迁移文件
migrate-create:
	@read -p "请输入迁移文件名: " name; \
	migrate create -ext sql -dir migrations $$name

# 重置数据库
db-reset:
	@echo "重置数据库..."
	@read -p "确定要重置数据库吗？这将删除所有数据 (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker-compose down -v; \
		docker-compose up -d postgres redis; \
		sleep 5; \
		make run; \
	else \
		echo "操作已取消"; \
	fi

# 备份数据库
db-backup:
	@echo "备份数据库..."
	@mkdir -p backups
	@docker exec exam_postgres pg_dump -U postgres online_exam_system > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "数据库备份完成"

# 恢复数据库
db-restore:
	@read -p "请输入备份文件路径: " file; \
	if [ -f "$$file" ]; then \
		docker exec -i exam_postgres psql -U postgres online_exam_system < "$$file"; \
		echo "数据库恢复完成"; \
	else \
		echo "备份文件不存在"; \
	fi