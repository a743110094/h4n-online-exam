# 更新日志 (Changelog)

本文档记录了项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- 开源项目文档和协议
- MIT 开源协议
- 贡献指南
- 行为准则
- 英文版 README

## [1.0.0] - 2024-01-27

### 新增
- 🎉 首次发布在线考试系统
- 👥 多角色用户系统（管理员、教师、学生）
- 📝 完整的考试管理功能
- 📚 题库管理系统
- 🤖 AI 问答助手集成
- 📊 成绩统计和分析
- 🏢 多租户架构支持
- ⚡ Redis 缓存优化
- 🐳 Docker 容器化部署
- 📱 响应式前端界面

### 技术特性
- Go + Gin 后端框架
- Vue 3 + TypeScript 前端
- PostgreSQL 数据库
- Redis 缓存系统
- JWT 身份认证
- GORM ORM 框架
- Element Plus UI 组件库
- Tailwind CSS 样式框架

### 性能优化
- 数据库索引优化
- 用户认证缓存
- 考试数据预热
- 压力测试验证（50VU，95%响应时间<100ms）

### 安全特性
- JWT Token 认证
- 密码 bcrypt 加密
- 多租户数据隔离
- CORS 跨域配置
- 权限中间件保护

### 部署支持
- Docker Compose 一键部署
- 环境变量配置
- 数据库迁移脚本
- 生产环境优化

## [0.9.0] - 2024-01-20

### 新增
- 压力测试脚本套件
- k6 性能测试工具集成
- 性能监控脚本
- 测试数据准备工具

### 优化
- 数据库查询性能
- 缓存策略改进
- API 响应时间优化

## [0.8.0] - 2024-01-15

### 新增
- Redis 缓存系统
- 用户认证缓存
- 考试数据缓存
- 缓存预热服务

### 修复
- 并发访问问题
- 内存泄漏修复
- 数据库连接池优化

## [0.7.0] - 2024-01-10

### 新增
- 数据库索引优化
- PostgreSQL 性能调优
- 查询优化

### 变更
- 数据库迁移脚本
- 索引策略调整

## [0.6.0] - 2024-01-05

### 新增
- AI 问答助手
- 聊天记录管理
- AI 接口集成

### 优化
- 用户体验改进
- 界面交互优化

## [0.5.0] - 2023-12-25

### 新增
- 统计分析功能
- 成绩报表生成
- 仪表板数据展示

### 修复
- 数据统计准确性
- 图表显示问题

## [0.4.0] - 2023-12-20

### 新增
- 考试管理系统
- 实时考试监控
- 答题功能
- 成绩计算

### 优化
- 考试流程优化
- 用户体验提升

## [0.3.0] - 2023-12-15

### 新增
- 试卷管理功能
- 自动组卷算法
- 手动组卷工具

### 变更
- 数据模型调整
- API 接口优化

## [0.2.0] - 2023-12-10

### 新增
- 题库管理系统
- 多种题型支持
- 题目分类管理

### 修复
- 数据验证问题
- 文件上传功能

## [0.1.0] - 2023-12-01

### 新增
- 基础项目架构
- 用户认证系统
- 角色权限管理
- 基础 CRUD 操作

### 技术栈
- Go + Gin 后端框架搭建
- Vue 3 前端框架搭建
- PostgreSQL 数据库集成
- JWT 认证实现

---

## 版本说明

### 版本号格式

本项目使用 [语义化版本](https://semver.org/lang/zh-CN/) 格式：`主版本号.次版本号.修订号`

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 变更类型

- `新增` - 新功能
- `变更` - 对现有功能的变更
- `弃用` - 即将移除的功能
- `移除` - 已移除的功能
- `修复` - 问题修复
- `安全` - 安全相关的修复

### 发布周期

- **主版本**：根据重大功能更新发布
- **次版本**：每月发布一次
- **修订版本**：根据 bug 修复需要发布

### 支持政策

- 当前主版本：完全支持
- 前一个主版本：安全更新支持 6 个月
- 更早版本：不再支持

---

[未发布]: https://github.com/your-username/online-exam-system/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-username/online-exam-system/releases/tag/v1.0.0
[0.9.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.9.0
[0.8.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.8.0
[0.7.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.7.0
[0.6.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.6.0
[0.5.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.5.0
[0.4.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.4.0
[0.3.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.3.0
[0.2.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.2.0
[0.1.0]: https://github.com/your-username/online-exam-system/releases/tag/v0.1.0