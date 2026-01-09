# 自测指南

本文档提供完整的自测步骤，帮助您验证代码实现是否满足需求。

## 📋 目录

1. [测试前准备](#测试前准备)
2. [环境配置检查](#环境配置检查)
3. [基础功能测试](#基础功能测试)
4. [Webhook 功能测试](#webhook-功能测试)
5. [统计 API 测试](#统计-api-测试)
6. [历史数据导入测试](#历史数据导入测试)
7. [数据库验证](#数据库验证)
8. [性能测试](#性能测试)
9. [问题排查](#问题排查)

---

## 测试前准备

### 1. 环境要求

- ✅ Go 1.21+ 已安装
- ✅ MySQL 5.7+ 或 PostgreSQL 12+ 已安装并运行
- ✅ 数据库已创建（默认：`gitlab_webhook`）
- ✅ 网络连接正常（用于测试 webhook）

### 2. 项目准备

```bash
# 1. 进入项目目录
cd gitlab_webhook_server

# 2. 安装依赖
go mod download

# 3. 复制环境变量文件
# Windows (PowerShell)
Copy-Item env.example .env

# Linux/Mac
cp env.example .env
```

### 3. 配置环境变量

编辑 `.env` 文件，配置以下关键项：

```env
# 服务器配置
PORT=3000
NODE_ENV=development
LOG_LEVEL=info

# Webhook 安全（可选，但建议配置）
GITLAB_WEBHOOK_SECRET=your_webhook_secret_here

# 数据库配置（根据实际情况修改）
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gitlab_webhook
DB_CHARSET=utf8mb4
DB_TIMEZONE=Asia/Shanghai

# 工作池配置（可选）
WORKER_POOL_WORKERS=10
WORKER_POOL_QUEUE_SIZE=1000

# 限流配置（可选）
RATE_LIMIT=100
RATE_LIMIT_WINDOW=1m

# GitLab API 配置（用于历史数据导入，可选）
GITLAB_BASE_URL=https://gitlab.com
GITLAB_TOKEN=your_gitlab_token_here
```

---

## 环境配置检查

### 1. 检查数据库连接

```bash
# MySQL
mysql -h localhost -u root -p -e "SHOW DATABASES;"

# PostgreSQL
psql -h localhost -U postgres -l
```

### 2. 启动服务

```bash
# 方式一：直接运行
go run cmd/server/main.go

# 方式二：使用 Makefile
make dev
```

**预期输出**：
```
🚀 服务器启动 port=3000 webhook_endpoint=http://localhost:3000/webhook health_endpoint=http://localhost:3000/health
数据库连接成功 type=mysql host=localhost database=gitlab_webhook
工作池已启动 workers=10 queue_size=1000
```

### 3. 验证服务启动

```bash
# 健康检查
curl http://localhost:3000/health
```

**预期响应**：
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## 基础功能测试

### 测试 1: 健康检查端点

```bash
curl http://localhost:3000/health
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含 `status: "ok"`
- ✅ 响应包含 `timestamp` 字段

### 测试 2: Webhook 测试端点

```bash
curl http://localhost:3000/webhook/test
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含 `message: "Webhook endpoint is ready"`

### 测试 3: 限流功能

```bash
# 快速发送多个请求（超过 RATE_LIMIT 限制）
for i in {1..110}; do curl http://localhost:3000/health; done
```

**预期结果**：
- ✅ 前 100 个请求成功（状态码 200）
- ✅ 超过限制的请求返回 429（Too Many Requests）
- ✅ 日志中记录限流信息

---

## Webhook 功能测试

### 测试 4: GitLab Webhook（自动检测）

**准备测试数据**：

创建 `test_payloads/gitlab_push.json`：

```json
{
  "object_kind": "push",
  "event_name": "push",
  "before": "95790bf891e76fee5e1747ab589903a6a1f80f22",
  "after": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
  "ref": "refs/heads/master",
  "ref_protected": true,
  "checkout_sha": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
  "message": "Hello World",
  "user_id": 4,
  "user_name": "John Smith",
  "user_username": "jsmith",
  "user_email": "john@example.com",
  "project_id": 15,
  "project": {
    "id": 15,
    "name": "Diaspora",
    "description": "",
    "web_url": "http://example.com/mike/diaspora",
    "path_with_namespace": "mike/diaspora",
    "namespace": "Mike",
    "visibility_level": 0,
    "default_branch": "master",
    "git_ssh_url": "git@example.com:mike/diaspora.git",
    "git_http_url": "http://example.com/mike/diaspora.git"
  },
  "commits": [
    {
      "id": "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "message": "Update Catalan translation to e38cb41.",
      "title": "Update Catalan translation to e38cb41.",
      "timestamp": "2011-12-12T14:27:31+02:00",
      "url": "http://example.com/mike/diaspora/commit/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "author": {
        "name": "Jordi Mallach",
        "email": "jordi@softcatala.org"
      },
      "added": ["CHANGELOG"],
      "modified": ["app/controller/application.rb"],
      "removed": []
    }
  ],
  "total_commits_count": 1,
  "repository": {
    "name": "Diaspora",
    "url": "git@example.com:mike/diaspora.git",
    "description": "",
    "homepage": "http://example.com/mike/diaspora",
    "git_http_url": "http://example.com/mike/diaspora.git",
    "git_ssh_url": "git@example.com:mike/diaspora.git",
    "visibility_level": 0
  }
}
```

**发送测试请求**：

```bash
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/json" \
  -H "X-Gitlab-Event: Push Hook" \
  -H "X-Gitlab-Token: your_webhook_secret_here" \
  -d @test_payloads/gitlab_push.json
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含 `status: "accepted"`
- ✅ 响应包含 `platform: "gitlab"`
- ✅ 响应包含 `event: "Push Hook"`
- ✅ 日志中记录 webhook 接收和处理信息
- ✅ 数据库中保存了提交记录

### 测试 5: Gitee Webhook

**准备测试数据**：

创建 `test_payloads/gitee_push.json`（格式类似 GitLab，但使用 Gitee 的字段名）：

```json
{
  "ref": "refs/heads/master",
  "before": "95790bf891e76fee5e1747ab589903a6a1f80f22",
  "after": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
  "total_commits_count": 1,
  "commits": [
    {
      "id": "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "message": "Update Catalan translation",
      "timestamp": "2011-12-12T14:27:31+02:00",
      "url": "http://example.com/mike/diaspora/commit/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "author": {
        "name": "Jordi Mallach",
        "email": "jordi@softcatala.org"
      },
      "added": ["CHANGELOG"],
      "modified": ["app/controller/application.rb"],
      "removed": []
    }
  ],
  "project": {
    "id": 15,
    "name": "Diaspora",
    "full_name": "mike/diaspora",
    "description": "",
    "html_url": "http://example.com/mike/diaspora",
    "default_branch": "master",
    "ssh_url": "git@example.com:mike/diaspora.git",
    "clone_url": "http://example.com/mike/diaspora.git"
  },
  "pusher": {
    "name": "John Smith",
    "email": "john@example.com"
  }
}
```

**发送测试请求**：

```bash
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/json" \
  -H "X-Gitee-Event: Push Hook" \
  -H "X-Gitee-Token: your_webhook_secret_here" \
  -d @test_payloads/gitee_push.json
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含 `platform: "gitee"`
- ✅ 数据库中保存了提交记录

### 测试 6: GitHub Webhook

**准备测试数据**：

创建 `test_payloads/github_push.json`：

```json
{
  "ref": "refs/heads/master",
  "before": "95790bf891e76fee5e1747ab589903a6a1f80f22",
  "after": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
  "commits": [
    {
      "id": "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "message": "Update Catalan translation",
      "timestamp": "2011-12-12T14:27:31+02:00",
      "url": "http://example.com/mike/diaspora/commit/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "author": {
        "name": "Jordi Mallach",
        "email": "jordi@softcatala.org"
      },
      "added": ["CHANGELOG"],
      "modified": ["app/controller/application.rb"],
      "removed": []
    }
  ],
  "repository": {
    "id": 15,
    "name": "Diaspora",
    "full_name": "mike/diaspora",
    "description": "",
    "html_url": "http://example.com/mike/diaspora",
    "default_branch": "master",
    "ssh_url": "git@example.com:mike/diaspora.git",
    "clone_url": "http://example.com/mike/diaspora.git",
    "private": false
  },
  "pusher": {
    "name": "John Smith",
    "email": "john@example.com"
  }
}
```

**发送测试请求**：

```bash
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -d @test_payloads/github_push.json
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含 `platform: "github"`
- ✅ 数据库中保存了提交记录

### 测试 7: Webhook Token 验证

**测试无效 token**：

```bash
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/json" \
  -H "X-Gitlab-Event: Push Hook" \
  -H "X-Gitlab-Token: invalid_token" \
  -d @test_payloads/gitlab_push.json
```

**预期结果**（如果配置了 `GITLAB_WEBHOOK_SECRET`）：
- ✅ HTTP 状态码：401（Unauthorized）
- ✅ 响应包含错误信息
- ✅ 日志中记录 token 验证失败

---

## 统计 API 测试

### 测试 8: 成员统计

```bash
# 查询成员统计（需要先有数据）
curl "http://localhost:3000/api/stats/member?email=jordi@softcatala.org&start_date=2011-01-01&end_date=2012-12-31"
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含：
  - `email`: 成员邮箱
  - `commit_count`: 提交数量
  - `total_added`: 总添加行数
  - `total_removed`: 总删除行数
  - `total_files`: 总变更文件数

**测试错误情况**：

```bash
# 缺少 email 参数
curl "http://localhost:3000/api/stats/member"
```

**预期结果**：
- ✅ HTTP 状态码：400（Bad Request）
- ✅ 响应包含错误信息

### 测试 9: 语言统计

```bash
curl "http://localhost:3000/api/stats/languages?email=jordi@softcatala.org&start_date=2011-01-01&end_date=2012-12-31"
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含语言统计数组，每个元素包含：
  - `language`: 编程语言
  - `total_added`: 该语言添加的行数
  - `total_removed`: 该语言删除的行数
  - `total_files`: 该语言的文件数

### 测试 10: 成员提交记录

```bash
curl "http://localhost:3000/api/stats/commits?email=jordi@softcatala.org&start_date=2011-01-01&end_date=2012-12-31"
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含：
  - `email`: 成员邮箱
  - `commits`: 提交记录数组
  - `count`: 提交数量

---

## 历史数据导入测试

### 测试 11: 导入项目提交记录

**前提条件**：
- ✅ 已配置 `GITLAB_BASE_URL` 和 `GITLAB_TOKEN`
- ✅ 有可访问的 GitLab 项目

```bash
curl -X POST http://localhost:3000/api/import/project \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "123",
    "since": "2024-01-01T00:00:00Z",
    "until": "2024-12-31T23:59:59Z",
    "batch_size": 100
  }'
```

**预期结果**：
- ✅ HTTP 状态码：202（Accepted）
- ✅ 响应包含：
  - `message`: "导入任务已启动"
  - `project_id`: 项目 ID
  - `status`: "processing"
- ✅ 日志中记录导入进度

### 测试 12: 查询导入状态

```bash
curl "http://localhost:3000/api/import/status?project_id=123"
```

**预期结果**：
- ✅ HTTP 状态码：200
- ✅ 响应包含导入状态信息

---

## 数据库验证

### 测试 13: 验证提交记录存储

**查询数据库**：

```sql
-- MySQL
USE gitlab_webhook;
SELECT * FROM commits ORDER BY created_at DESC LIMIT 10;

-- PostgreSQL
\c gitlab_webhook;
SELECT * FROM commits ORDER BY created_at DESC LIMIT 10;
```

**验证点**：
- ✅ 提交记录已保存
- ✅ `commit_id` 和 `project_id` 组合唯一
- ✅ 时间戳正确
- ✅ 作者信息完整
- ✅ 项目信息完整

### 测试 14: 验证文件变更记录

```sql
SELECT * FROM commit_files ORDER BY created_at DESC LIMIT 10;
```

**验证点**：
- ✅ 文件变更记录已保存
- ✅ `change_type` 正确（added/modified/removed）
- ✅ 文件路径和扩展名正确
- ✅ 语言检测正确

### 测试 15: 验证语言统计

```sql
SELECT * FROM commit_languages ORDER BY created_at DESC LIMIT 10;
```

**验证点**：
- ✅ 语言统计已保存
- ✅ 语言名称正确
- ✅ 行数统计正确

---

## 性能测试

### 测试 16: 并发 Webhook 处理

```bash
# 使用 Apache Bench 或类似工具
ab -n 1000 -c 10 -p test_payloads/gitlab_push.json -T application/json \
   -H "X-Gitlab-Event: Push Hook" \
   http://localhost:3000/webhook
```

**验证点**：
- ✅ 所有请求都返回 200 或 202
- ✅ 没有错误日志
- ✅ 数据库记录正确
- ✅ 工作池正常工作

### 测试 17: 限流功能

```bash
# 快速发送大量请求
for i in {1..200}; do
  curl -s http://localhost:3000/health > /dev/null &
done
wait
```

**验证点**：
- ✅ 前 100 个请求成功
- ✅ 超过限制的请求返回 429
- ✅ 日志中记录限流信息

---

## 问题排查

### 常见问题

#### 1. 数据库连接失败

**症状**：
```
数据库初始化失败: connection refused
```

**解决方案**：
- ✅ 检查数据库服务是否运行
- ✅ 检查数据库配置（host, port, user, password）
- ✅ 检查数据库是否已创建
- ✅ 检查防火墙设置

#### 2. Webhook 接收但未保存

**症状**：
- Webhook 返回 200，但数据库中没有记录

**排查步骤**：
1. 检查日志中的错误信息
2. 检查工作池是否正常运行
3. 检查数据库连接是否正常
4. 检查提交记录是否已存在（唯一性检查）

#### 3. Token 验证失败

**症状**：
```
Webhook token 验证失败
```

**解决方案**：
- ✅ 检查 `.env` 中的 `GITLAB_WEBHOOK_SECRET` 配置
- ✅ 确保 GitLab/Gitee 中配置的 token 与 `.env` 中的一致
- ✅ 如果不需要验证，可以留空 `GITLAB_WEBHOOK_SECRET`

#### 4. 统计 API 返回空数据

**症状**：
- API 返回成功，但数据为空

**排查步骤**：
1. 检查数据库中是否有数据
2. 检查 email 参数是否正确
3. 检查时间范围是否合理
4. 检查日志中的查询信息

### 日志查看

**查看服务日志**：
```bash
# 如果使用 systemd
journalctl -u gitlab-webhook-server -f

# 如果直接运行
# 日志会输出到控制台
```

**关键日志信息**：
- `🚀 服务器启动` - 服务启动成功
- `数据库连接成功` - 数据库连接成功
- `收到 Webhook 事件` - Webhook 接收
- `任务执行成功` - 任务处理成功
- `记录提交失败` - 提交保存失败（需要关注）

---

## 测试检查清单

### 基础功能 ✅

- [ ] 服务正常启动
- [ ] 健康检查端点正常
- [ ] Webhook 测试端点正常
- [ ] 限流功能正常

### Webhook 功能 ✅

- [ ] GitLab webhook 接收和处理
- [ ] Gitee webhook 接收和处理
- [ ] GitHub webhook 接收和处理
- [ ] Token 验证功能（如果配置）
- [ ] 平台自动检测功能

### 数据存储 ✅

- [ ] 提交记录正确保存
- [ ] 文件变更记录正确保存
- [ ] 语言统计正确保存
- [ ] 唯一性检查正常工作

### 统计 API ✅

- [ ] 成员统计 API 正常
- [ ] 语言统计 API 正常
- [ ] 提交记录查询 API 正常
- [ ] 参数验证正常

### 历史数据导入 ✅

- [ ] 导入任务启动正常
- [ ] 导入状态查询正常（如果实现）

### 性能 ✅

- [ ] 并发处理正常
- [ ] 限流功能正常
- [ ] 工作池正常工作

---

## 测试完成标准

### 必须通过 ✅

1. ✅ 服务正常启动
2. ✅ 至少一个平台的 webhook 可以正常接收和处理
3. ✅ 数据正确保存到数据库
4. ✅ 统计 API 可以正常查询数据

### 建议通过 ⚠️

1. ⚠️ 所有三个平台的 webhook 都可以正常工作
2. ⚠️ Token 验证功能正常
3. ⚠️ 历史数据导入功能正常（如果配置了 GitLab API）

---

## 下一步

测试通过后，您可以：

1. **配置生产环境**：参考 [Docker 部署指南](./docker_deployment.md)
2. **配置 GitLab/Gitee/GitHub Webhook**：在代码仓库中配置 webhook URL
3. **监控服务**：设置日志监控和告警
4. **性能优化**：根据实际负载调整工作池和限流配置

---

## 获取帮助

如果遇到问题：

1. 查看 [代码审查报告](./code_review_report.md)
2. 查看 [开发文档](./DEVELOPMENT.md)
3. 检查日志文件
4. 查看项目 Issue

---

**祝测试顺利！** 🎉
