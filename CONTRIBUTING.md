# 贡献指南 (Contributing Guide)

感谢您对在线考试系统项目的关注！我们欢迎所有形式的贡献，包括但不限于：

- 🐛 Bug 报告
- 💡 功能建议
- 📝 文档改进
- 🔧 代码贡献
- 🌐 翻译工作

## 🚀 快速开始

### 环境准备

1. **Fork 项目**
   - 点击页面右上角的 "Fork" 按钮
   - 克隆你的 fork 到本地

```bash
git clone https://github.com/your-username/online-exam-system.git
cd online-exam-system
```

2. **设置开发环境**

```bash
# 添加上游仓库
git remote add upstream https://github.com/original-owner/online-exam-system.git

# 安装依赖
# 后端
cd backend && go mod download

# 前端
cd frontend && npm install
```

3. **启动开发环境**

```bash
# 启动数据库
docker-compose up -d postgres redis

# 启动后端 (终端1)
cd backend && go run main.go

# 启动前端 (终端2)
cd frontend && npm run dev
```

## 📋 开发流程

### 1. 创建分支

```bash
# 同步最新代码
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature-name
# 或者修复分支
git checkout -b fix/your-bug-fix
```

### 2. 开发规范

#### 代码风格

**Go 代码规范:**
- 使用 `gofmt` 格式化代码
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 函数和变量命名使用驼峰命名法
- 添加必要的注释，特别是公开的函数和结构体

**Vue/TypeScript 代码规范:**
- 使用 ESLint 和 Prettier 格式化代码
- 组件名使用 PascalCase
- 文件名使用 kebab-case
- 使用 TypeScript 类型注解

#### 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型说明:**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式化
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例:**
```
feat(auth): add JWT token refresh mechanism

fix(exam): resolve timer synchronization issue

docs: update API documentation for exam endpoints
```

### 3. 测试

在提交代码前，请确保：

```bash
# 后端测试
cd backend
go test ./...
go vet ./...

# 前端测试
cd frontend
npm run lint
npm run test:unit
npm run build
```

### 4. 提交 Pull Request

1. **推送分支**
```bash
git push origin feature/your-feature-name
```

2. **创建 Pull Request**
   - 访问 GitHub 页面
   - 点击 "New Pull Request"
   - 填写详细的描述

3. **PR 模板**
```markdown
## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 性能优化
- [ ] 其他

## 变更描述
简要描述你的变更内容...

## 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手动测试通过

## 截图 (如适用)

## 相关 Issue
Closes #123
```

## 🐛 Bug 报告

发现 Bug？请创建 Issue 并包含以下信息：

### Bug 报告模板

```markdown
## Bug 描述
简要描述遇到的问题...

## 复现步骤
1. 访问 '...'
2. 点击 '...'
3. 滚动到 '...'
4. 看到错误

## 期望行为
描述你期望发生的情况...

## 实际行为
描述实际发生的情况...

## 环境信息
- OS: [e.g. macOS 12.0]
- Browser: [e.g. Chrome 95.0]
- Go Version: [e.g. 1.19]
- Node Version: [e.g. 18.0]

## 附加信息
添加任何其他相关信息、截图或日志...
```

## 💡 功能建议

有新的想法？我们很乐意听到！请创建 Issue 并使用以下模板：

### 功能建议模板

```markdown
## 功能描述
简要描述你建议的功能...

## 问题背景
这个功能解决了什么问题？

## 解决方案
描述你希望的解决方案...

## 替代方案
描述你考虑过的其他替代方案...

## 附加信息
添加任何其他相关信息或截图...
```

## 📚 文档贡献

文档同样重要！你可以帮助我们：

- 修复文档中的错误
- 改进现有文档的清晰度
- 添加缺失的文档
- 翻译文档到其他语言

## 🌐 国际化

我们欢迎翻译贡献！目前支持的语言：

- 中文 (zh-CN)
- English (en-US)

如果你想添加新的语言支持，请：

1. 在 `frontend/src/locales/` 目录下添加新的语言文件
2. 更新语言选择器组件
3. 测试所有页面的翻译

## 🔍 代码审查

所有的 Pull Request 都需要经过代码审查。审查者会关注：

- 代码质量和可读性
- 测试覆盖率
- 性能影响
- 安全性
- 文档完整性

## 📞 获取帮助

如果你在贡献过程中遇到问题，可以通过以下方式获取帮助：

- 创建 [Discussion](https://github.com/your-username/online-exam-system/discussions)
- 在相关 Issue 中评论
- 联系维护者

## 🎉 贡献者认可

我们重视每一个贡献！所有贡献者都会被列在项目的贡献者列表中。

感谢你的贡献！🙏