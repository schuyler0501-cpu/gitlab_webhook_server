# 🚀 快速启动指南

## 第一步：安装 Go

确保已安装 Go 1.21 或更高版本：

```bash
go version
```

如果未安装，请访问 [Go 官网](https://go.dev/dl/) 下载安装。

## 第二步：安装依赖

```bash
go mod download
go mod tidy
```

或使用 Makefile:

```bash
make deps
```

## 第三步：配置环境变量

复制环境变量模板文件：

```bash
# Windows (PowerShell)
Copy-Item env.example .env

# Linux/Mac
cp env.example .env
```

编辑 `.env` 文件，至少配置以下内容：

```env
PORT=3000
GITLAB_WEBHOOK_SECRET=your_secret_token_here
```

## 第四步：启动开发服务器

### 方式一：使用 Air（推荐，支持热重载）

首先安装 Air:
```bash
go install github.com/air-verse/air@latest
```

然后运行:
```bash
make dev
```

### 方式二：直接运行

```bash
go run cmd/server/main.go
```

看到以下输出表示启动成功：

```
🚀 服务器启动在端口 3000
📡 Webhook 端点: http://localhost:3000/webhook
💚 健康检查: http://localhost:3000/health
```

## 第五步：测试服务

在浏览器中访问：

- 健康检查: http://localhost:3000/health
- Webhook 测试: http://localhost:3000/webhook/test

或使用 curl:

```bash
curl http://localhost:3000/health
curl http://localhost:3000/webhook/test
```

## 第六步：配置 GitLab Webhook

1. 登录 GitLab，进入你的项目
2. 进入 **Settings** → **Webhooks**
3. 填写以下信息：
   - **URL**: `http://your-server-ip:3000/webhook`
   - **Secret token**: 与 `.env` 文件中的 `GITLAB_WEBHOOK_SECRET` 保持一致
   - **Trigger**: 勾选 `Push events`
4. 点击 **Add webhook**

## 验证 Webhook

在 GitLab 项目中提交代码，然后查看服务器日志，应该能看到类似以下输出：

```
收到 Webhook 事件: Push Hook
📝 记录代码提交: { commit_id: '...', author: '...', ... }
📊 提交统计: { added_files: 2, modified_files: 1, ... }
```

## 常用命令

```bash
# 开发模式（热重载，需要安装 air）
make dev

# 构建项目
make build

# 运行生产版本
make run

# 代码检查（需要安装 golangci-lint）
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

## 安装推荐工具

### Air（热重载）
```bash
go install github.com/air-verse/air@latest
```

### golangci-lint（代码检查）
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### goimports（自动导入管理）
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

## 遇到问题？

1. **端口被占用**: 修改 `.env` 文件中的 `PORT` 值
2. **依赖下载失败**: 检查网络连接，或配置 Go 代理：
   ```bash
   go env -w GOPROXY=https://goproxy.cn,direct
   ```
3. **编译错误**: 运行 `go mod tidy` 整理依赖
4. **找不到命令**: 确保 `$GOPATH/bin` 或 `$HOME/go/bin` 在 PATH 中

## 下一步

- 阅读 [开发规范文档](./docs/DEVELOPMENT.md) 了解项目结构
- 查看 [AI 辅助开发指南](./docs/AI_CODING_GUIDE.md) 学习如何与 AI 协作
- 开始实现你的功能！

