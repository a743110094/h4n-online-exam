# 在线考试系统 (Online Exam System)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-green.svg)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

一个功能完整的在线考试系统，支持多租户、AI助手、实时监控等特性。基于现代化技术栈构建，提供高性能、高可用的考试解决方案。

[English](./README_EN.md) | 中文

## ✨ 功能特性

### 🎯 核心功能
- **多角色支持**: 管理员、教师、学生三种角色，权限分离
- **考试管理**: 完整的考试流程，从创建到统计分析
- **题库管理**: 支持单选、多选、判断、简答等多种题型
- **智能组卷**: 手动组卷和AI自动组卷
- **实时监控**: 考试过程实时监控，防作弊机制
- **成绩统计**: 详细的成绩分析和统计报表

### 🚀 技术特性
- **多租户架构**: 支持多机构独立使用
- **Redis缓存**: 高性能缓存，支持高并发
- **AI集成**: 内置AI问答助手
- **响应式设计**: 支持PC和移动端
- **Docker部署**: 一键部署，开箱即用
- **性能优化**: 数据库索引优化，压力测试验证

## 🛠️ 技术栈

### 后端
- **框架**: Go + Gin
- **数据库**: PostgreSQL + Redis
- **ORM**: GORM
- **认证**: JWT
- **缓存**: Redis
- **容器**: Docker + Docker Compose

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **样式**: Tailwind CSS

## 📦 快速开始

### 环境要求

- Docker & Docker Compose
- Go 1.19+ (开发环境)
- Node.js 18+ (开发环境)

### 一键部署 (推荐)

```bash
# 克隆项目
git clone https://github.com/your-username/online-exam-system.git
cd online-exam-system

# 启动所有服务
docker-compose up -d

# 等待服务启动完成后访问
# 前端: http://localhost:5173
# 后端API: http://localhost:8080
```

### 开发环境部署

#### 1. 启动数据库服务

```bash
# 启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis
```

#### 2. 后端服务

```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件配置数据库连接

# 运行服务
go run main.go
```

#### 3. 前端服务

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 🎮 默认账号

系统初始化后提供以下测试账号（密码均为 `admin123`）：

| 角色 | 用户名 | 密码 | 说明 |
|------|--------|------|------|
| 管理员 | admin | admin123 | 系统管理员 |
| 教师 | teacher1 | admin123 | 教师账号 |
| 学生 | student1 | admin123 | 学生账号 |

**⚠️ 生产环境请务必修改默认密码！**

## 📖 API 文档

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/profile` - 获取用户信息
- `PUT /api/v1/auth/profile` - 更新用户信息

### 考试管理
- `GET /api/v1/exams` - 获取考试列表
- `POST /api/v1/exams` - 创建考试
- `POST /api/v1/exams/:id/start` - 开始考试
- `POST /api/v1/exams/:id/submit` - 提交考试

### 题目管理
- `GET /api/v1/questions` - 获取题目列表
- `POST /api/v1/questions` - 创建题目
- `PUT /api/v1/questions/:id` - 更新题目
- `DELETE /api/v1/questions/:id` - 删除题目

更多API文档请参考 [API文档](./docs/API.md)

## 🏗️ 项目结构

```
.
├── backend/                 # 后端服务
│   ├── controllers/         # 控制器
│   ├── models/             # 数据模型
│   ├── middleware/         # 中间件
│   ├── services/           # 业务服务
│   ├── cache/              # 缓存服务
│   └── main.go             # 入口文件
├── frontend/               # 前端应用
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── views/          # 页面
│   │   ├── stores/         # 状态管理
│   │   └── api/            # API接口
│   └── package.json
├── load-test/              # 压力测试脚本
├── docker-compose.yml      # Docker编排
└── README.md
```

## 🧪 测试

### 压力测试

项目包含完整的k6压力测试脚本，支持50VU并发，95%响应时间<100ms：

```bash
# 安装k6
brew install k6  # macOS
# 或 sudo apt install k6  # Ubuntu

# 运行压力测试
cd load-test
./run-load-tests.sh
```

### 单元测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test:unit
```

## 📊 性能优化

- **数据库索引**: 针对高频查询字段添加索引
- **Redis缓存**: 缓存热点数据，减少数据库压力
- **连接池**: 数据库连接池优化
- **压缩**: Gzip压缩减少传输大小
- **CDN**: 静态资源CDN加速

## 🚀 部署指南

### 生产环境部署

1. **环境配置**
   ```bash
   # 设置生产环境变量
   export GIN_MODE=release
   export DB_HOST=your-db-host
   export REDIS_HOST=your-redis-host
   ```

2. **数据库迁移**
   ```bash
   # 执行数据库迁移
   psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f backend/migrations/add_indexes.sql
   ```

3. **启动服务**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### 监控和日志

- 使用 Prometheus + Grafana 监控系统性能
- 集成 ELK Stack 进行日志分析
- 配置告警规则，及时发现问题

## 🤝 贡献指南

我们欢迎所有形式的贡献！请查看 [贡献指南](./CONTRIBUTING.md) 了解详情。

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](./LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系我们

- 项目主页: [GitHub](https://github.com/your-username/online-exam-system)
- 问题反馈: [Issues](https://github.com/your-username/online-exam-system/issues)
- 讨论交流: [Discussions](https://github.com/your-username/online-exam-system/discussions)

## 🌟 Star History

如果这个项目对你有帮助，请给我们一个 ⭐️！

[![Star History Chart](https://api.star-history.com/svg?repos=your-username/online-exam-system&type=Date)](https://star-history.com/#your-username/online-exam-system&Date)