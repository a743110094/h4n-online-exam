# 在线考试系统后端

基于 Go + Gin + GORM + PostgreSQL 构建的在线考试系统后端服务。

## 功能特性

### 用户管理
- 用户注册/登录
- 角色管理（管理员、教师、学生）
- 用户信息管理
- 密码修改/重置

### 考试管理
- 科目管理
- 题目管理（单选、多选、判断、简答）
- 试卷管理（手动组卷、自动组卷）
- 考试管理（创建、开始、结束）
- 答题功能
- 成绩统计

### AI 功能
- AI 问答助手
- 聊天记录管理

### 统计分析
- 学生成绩统计
- 教师教学统计
- 管理员仪表板
- 考试分析报告

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT
- **密码加密**: bcrypt
- **配置管理**: godotenv

## 项目结构

```
backend/
├── config/          # 配置管理
├── controllers/     # 控制器
├── database/        # 数据库连接和迁移
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由定义
├── scripts/         # 数据库脚本
├── uploads/         # 文件上传目录
├── main.go          # 程序入口
├── go.mod           # Go模块文件
├── go.sum           # 依赖校验文件
├── Makefile         # 构建脚本
├── docker-compose.yml # Docker编排文件
└── README.md        # 项目说明
```

## 快速开始

### 环境要求

- Go 1.19+
- PostgreSQL 13+
- Redis 6+

### 安装依赖

```bash
go mod download
```

### 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息
```

### 数据库初始化

```bash
# 使用 Docker Compose 启动数据库服务
docker-compose up -d postgres redis

# 或者手动创建数据库并导入初始数据
psql -U postgres -c "CREATE DATABASE online_exam_system;"
psql -U postgres -d online_exam_system -f scripts/init_db.sql
```

### 运行项目

```bash
# 开发模式
make dev

# 或者直接运行
go run main.go
```

### 使用 Docker

```bash
# 启动完整环境
make docker-up

# 停止服务
make docker-down
```

## API 文档

### 认证相关

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册（管理员）

### 用户管理

- `GET /api/v1/user/profile` - 获取用户信息
- `PUT /api/v1/user/profile` - 更新用户信息
- `PUT /api/v1/user/password` - 修改密码

### 管理员接口

- `GET /api/v1/admin/users` - 获取用户列表
- `POST /api/v1/admin/users` - 创建用户
- `PUT /api/v1/admin/users/:id` - 更新用户
- `DELETE /api/v1/admin/users/:id` - 删除用户
- `GET /api/v1/admin/dashboard` - 仪表板统计

### 教师接口

- `GET /api/v1/questions` - 获取题目列表
- `POST /api/v1/teacher/questions` - 创建题目
- `PUT /api/v1/teacher/questions/:id` - 更新题目
- `DELETE /api/v1/teacher/questions/:id` - 删除题目
- `POST /api/v1/teacher/papers` - 创建试卷
- `POST /api/v1/teacher/papers/auto` - 自动组卷
- `POST /api/v1/teacher/exams` - 创建考试

### 学生接口

- `GET /api/v1/exams/student` - 获取学生考试列表
- `POST /api/v1/exams/:id/start` - 开始考试
- `POST /api/v1/exams/:exam_id/answers` - 提交答案
- `POST /api/v1/exams/:exam_id/answers/submit` - 提交试卷
- `GET /api/v1/stats/student` - 学生统计

### AI 接口

- `POST /api/v1/ai/chat` - AI 问答
- `GET /api/v1/ai/history` - 获取聊天记录
- `DELETE /api/v1/ai/history` - 清空聊天记录

## 默认账号

系统初始化后会创建以下默认账号（密码均为 `admin123`）：

- **管理员**: `admin`
- **教师**: `teacher1`, `teacher2`
- **学生**: `student1`, `student2`, `student3`

**注意**: 生产环境部署前请务必修改默认密码！

## 开发指南

### 添加新的控制器

1. 在 `controllers/` 目录下创建新的控制器文件
2. 在 `routes/routes.go` 中注册路由
3. 如需要权限控制，使用相应的中间件

### 添加新的数据模型

1. 在 `models/models.go` 中定义新的结构体
2. 在 `database/database.go` 的 `AutoMigrate()` 函数中添加新模型

### 中间件使用

- `AuthMiddleware()` - 身份认证
- `RoleMiddleware(roles...)` - 角色权限控制

## 部署

### 生产环境配置

1. 设置环境变量 `GIN_MODE=release`
2. 配置生产数据库连接
3. 设置强密码的 JWT Secret
4. 配置 HTTPS
5. 设置日志级别

### Docker 部署

```bash
# 构建镜像
make docker-build

# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d
```

## 常用命令

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 构建二进制文件
make build

# 清理构建文件
make clean

# 数据库迁移
make migrate

# 重置数据库
make reset-db
```

## 许可证

MIT License