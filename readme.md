# GitLab Webhook Server

## 📋 项目介绍

本项目用于作为 GitLab 的 webhook 服务，目的是为了获取开发团队成员代码提交记录，记录代码提交的情况，统计成员的代码。

核心目的是在效能度量的应用场景使用，统计团队成员在领取了开发任务后，实际的代码产出与任务估时的情况展示，确认大家的饱和度。

## 🚀 快速开始

### 前置要求

- Go >= 1.21（本地开发）
- Make (可选，用于使用 Makefile 命令)
- Docker >= 20.10 和 Docker Compose >= 2.0（Docker 部署）

### 方式一：Docker 部署（推荐）

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd gitlab-webhook-server
   ```

2. **配置环境变量**

   ```bash
   # Windows (PowerShell)
   Copy-Item env.example .env
   
   # Linux/Mac
   cp env.example .env
   ```

   然后编辑 `.env` 文件，填入你的配置信息。

3. **启动服务**

   ```bash
   # 构建并启动所有服务（包括数据库）
   docker-compose up -d
   
   # 查看日志
   docker-compose logs -f
   
   # 查看服务状态
   docker-compose ps
   ```

4. **访问服务**
   - 服务器地址: <http://localhost:3000>
   - Webhook 端点: <http://localhost:3000/webhook>
   - 健康检查: <http://localhost:3000/health>

**详细说明**: 请参考 [Docker 部署指南](./docs/docker_deployment.md)

### 方式二：本地开发

1. **安装依赖**

   ```bash
   go mod download
   go mod tidy
   ```

   或使用 Makefile:

   ```bash
   make deps
   ```

2. **配置环境变量**

   ```bash
   # Windows (PowerShell)
   Copy-Item env.example .env
   
   # Linux/Mac
   cp env.example .env
   ```

   然后编辑 `.env` 文件，填入你的配置信息。

3. **启动开发服务器**

   ```bash
   # 使用 Makefile（推荐）
   make dev
   
   # 或直接运行
   go run cmd/server/main.go
   ```

4. **访问服务**
   - 服务器地址: <http://localhost:3000>
   - Webhook 端点: <http://localhost:3000/webhook>
   - 健康检查: <http://localhost:3000/health>

### 在 GitLab 中配置 Webhook

1. 进入你的 GitLab 项目
2. 进入 **Settings** → **Webhooks**
3. 添加新的 Webhook:
   - URL: `http://your-server:3000/webhook`
   - Secret token: 与 `.env` 中的 `GITLAB_WEBHOOK_SECRET` 保持一致
   - 触发事件: 选择 `Push events` 和其他需要的事件

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **日志库**: Uber Zap
- **配置管理**: godotenv
- **代码质量**: golangci-lint
- **测试**: Go 标准测试框架

## 📁 项目结构

```
gitlab_webhook_server/
├── cmd/                    # 应用入口
│   └── server/            # 服务器主程序
│       └── main.go
├── internal/              # 内部包（不对外暴露）
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── model/            # 数据模型
│   ├── router/           # 路由定义
│   ├── service/          # 业务逻辑层
│   │   └── commit/       # 提交相关服务
│   └── logger/           # 日志工具
├── docs/                 # 文档
│   ├── DEVELOPMENT.md    # 开发规范
│   └── AI_CODING_GUIDE.md # AI 辅助开发指南
├── bin/                  # 编译输出（自动生成）
├── go.mod                # Go 模块定义
└── Makefile             # 构建脚本
```

## 📝 开发命令

```bash
# 开发模式（热重载，需要安装 air）
make dev

# 构建项目
make build

# 运行生产版本
make run

# 代码检查
make lint

# 格式化代码
make fmt

# 运行测试
make test

# 清理构建文件
make clean

# 安装依赖
make deps

# 查看所有命令
make help
```

## 📚 文档

- [Docker 部署指南](./docs/docker_deployment.md) - Docker 和 Docker Compose 部署详细说明
- [开发规范](./docs/DEVELOPMENT.md) - 详细的开发规范和最佳实践
- [AI 辅助开发指南](./docs/AI_CODING_GUIDE.md) - 如何更好地与 AI 协作开发
- [快速启动指南](./QUICKSTART.md) - 快速上手指南
- [数据库设计文档](./docs/database_design.md) - 数据库表结构设计
- [数据库配置指南](./docs/database_setup.md) - 数据库配置和迁移
- [性能优化文档](./docs/performance_optimization.md) - 并发优化和历史数据导入
- [Cursor AI 配置说明](./README_CURSOR.md) - Cursor AI 配置和使用指南
- [有效对话记录系统](./README_CONVERSATIONS.md) - 有效对话自动记录机制

## 🤖 AI 辅助开发 (Vibe Coding)

本项目已配置完善的 **Cursor AI / Vibe Coding** 开发环境，让 AI 能够更好地理解和协助开发：

### ✨ Cursor AI 配置

- ✅ **`.cursorrules`** - 完整的项目编码规范和最佳实践
- ✅ **`.cursor/context.md`** - 项目上下文和业务逻辑说明
- ✅ **`.cursor/instructions.md`** - AI 代码生成指令和模板
- ✅ **`.cursor/prompts.md`** - 常用提示词模板库
- ✅ **`.cursorignore`** - 优化 AI 索引，提高响应速度

### 🎯 核心特性

- ✅ Go 标准项目结构，清晰的代码组织
- ✅ 代码质量工具（golangci-lint）
- ✅ 详细的开发文档和示例代码
- ✅ 规范的错误处理和日志记录
- ✅ 完善的类型定义和接口设计
- ✅ **AI 自动遵循项目规范** - 生成的代码符合项目标准

### 📚 相关文档

- [Cursor AI 配置说明](./README_CURSOR.md) - 详细的 Cursor 配置和使用指南
- [AI 辅助开发指南](./docs/AI_CODING_GUIDE.md) - 如何更好地与 AI 协作
- [开发规范](./docs/DEVELOPMENT.md) - 详细的开发规范和最佳实践

### 🚀 快速开始使用 AI

1. **启用 Vibe Coding**: 在 Cursor 设置中启用 Vibe Coding 功能
2. **直接对话**: 与 AI 对话时，AI 会自动遵循 `.cursorrules` 中的规范
3. **使用模板**: 参考 `.cursor/prompts.md` 中的提示词模板

**示例**:

```
用户: 帮我添加一个查询提交记录的 API 端点

AI: 会自动在正确的目录创建文件，遵循命名规范，
    添加错误处理和日志记录，符合项目架构模式。
```

## 🛠️ 推荐工具

- **Air**: 热重载工具

  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **golangci-lint**: 代码检查工具

  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **goimports**: 自动导入管理

  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  ```

## 📄 License

MIT
