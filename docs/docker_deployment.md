# Docker 部署指南

本文档介绍如何使用 Docker 和 Docker Compose 部署 GitLab Webhook Server。

## 📋 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [生产环境部署](#生产环境部署)
- [数据库迁移](#数据库迁移)
- [健康检查](#健康检查)
- [故障排查](#故障排查)

## 前置要求

- Docker >= 20.10
- Docker Compose >= 2.0

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd gitlab-webhook-server
```

### 2. 配置环境变量

复制环境变量示例文件：

```bash
cp env.example .env
```

编辑 `.env` 文件，配置必要的环境变量：

```bash
# 服务器配置
PORT=3000
NODE_ENV=production

# GitLab Webhook 配置
GITLAB_WEBHOOK_SECRET=your_webhook_secret_here

# 数据库配置
DB_TYPE=mysql
DB_HOST=mysql  # Docker Compose 中使用服务名
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_secure_password
DB_NAME=gitlab_webhook
DB_CHARSET=utf8mb4
DB_TIMEZONE=Asia/Shanghai

# GitLab API 配置（可选，用于历史数据导入）
GITLAB_BASE_URL=https://gitlab.com
GITLAB_TOKEN=your_gitlab_token_here
```

### 3. 启动服务

使用 Docker Compose 启动所有服务：

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 4. 验证部署

```bash
# 健康检查
curl http://localhost:3000/health

# 查看服务日志
docker-compose logs webhook-server
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| `PORT` | 服务端口 | 3000 | 否 |
| `NODE_ENV` | 运行环境 | production | 否 |
| `LOG_LEVEL` | 日志级别 | info | 否 |
| `GITLAB_WEBHOOK_SECRET` | Webhook 密钥 | - | 是 |
| `DB_TYPE` | 数据库类型 | mysql | 否 |
| `DB_HOST` | 数据库主机 | mysql | 是 |
| `DB_PORT` | 数据库端口 | 3306 | 否 |
| `DB_USER` | 数据库用户 | root | 是 |
| `DB_PASSWORD` | 数据库密码 | - | 是 |
| `DB_NAME` | 数据库名称 | gitlab_webhook | 是 |
| `WORKER_POOL_WORKERS` | 工作池协程数 | 10 | 否 |
| `WORKER_POOL_QUEUE_SIZE` | 工作池队列大小 | 1000 | 否 |
| `RATE_LIMIT` | 限流数量 | 100 | 否 |
| `RATE_LIMIT_WINDOW` | 限流时间窗口 | 1m | 否 |
| `GITLAB_BASE_URL` | GitLab 地址 | https://gitlab.com | 否 |
| `GITLAB_TOKEN` | GitLab Token | - | 否 |

### 数据库选择

#### 使用 MySQL（默认）

```yaml
# docker-compose.yml 中已配置 MySQL
services:
  mysql:
    image: mysql:8.0
    # ...
```

#### 使用 PostgreSQL

1. 注释掉 `docker-compose.yml` 中的 MySQL 服务
2. 取消注释 PostgreSQL 服务
3. 修改环境变量：

```bash
DB_TYPE=postgresql
DB_HOST=postgres
DB_PORT=5432
```

## 生产环境部署

### 使用生产环境配置

```bash
# 使用生产环境配置启动
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

生产环境配置包括：
- 资源限制（CPU、内存）
- 日志轮转配置
- 数据库性能优化

### 使用反向代理

生产环境建议使用 Nginx 或 Traefik 作为反向代理：

```nginx
# nginx.conf 示例
server {
    listen 80;
    server_name webhook.example.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 使用外部数据库

如果使用外部数据库（如云数据库），修改环境变量：

```bash
# .env
DB_HOST=your-database-host.com
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
```

并在 `docker-compose.yml` 中移除数据库服务，移除 `depends_on` 配置。

## 数据库迁移

### 自动迁移

应用启动时会自动执行数据库迁移（通过 GORM AutoMigrate）。

### 手动迁移

如果需要手动执行 SQL 迁移文件：

```bash
# 进入 MySQL 容器
docker-compose exec mysql bash

# 执行迁移
mysql -u root -p gitlab_webhook < /docker-entrypoint-initdb.d/001_create_tables_mysql.sql
mysql -u root -p gitlab_webhook < /docker-entrypoint-initdb.d/002_optimize_tables_mysql.sql
mysql -u root -p gitlab_webhook < /docker-entrypoint-initdb.d/003_add_webhook_fields_mysql.sql
```

## 健康检查

### 容器健康检查

Docker Compose 配置了健康检查：

```bash
# 查看健康状态
docker-compose ps

# 查看健康检查日志
docker inspect gitlab-webhook-server | grep -A 10 Health
```

### 应用健康检查

应用提供 `/health` 端点：

```bash
curl http://localhost:3000/health
```

## 故障排查

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs webhook-server
docker-compose logs mysql

# 实时跟踪日志
docker-compose logs -f webhook-server
```

### 常见问题

#### 1. 数据库连接失败

**问题**: 应用无法连接到数据库

**解决方案**:
- 检查数据库服务是否启动: `docker-compose ps`
- 检查环境变量配置是否正确
- 检查数据库健康状态: `docker-compose exec mysql mysqladmin ping -h localhost -u root -p`

#### 2. 端口冲突

**问题**: 端口 3000 已被占用

**解决方案**:
- 修改 `.env` 中的 `PORT` 变量
- 或修改 `docker-compose.yml` 中的端口映射: `"8080:3000"`

#### 3. 数据库迁移失败

**问题**: 数据库表未创建

**解决方案**:
- 检查数据库连接配置
- 查看应用日志: `docker-compose logs webhook-server`
- 手动执行迁移（参考上面的"手动迁移"部分）

#### 4. Webhook 无法接收请求

**问题**: GitLab Webhook 请求失败

**解决方案**:
- 检查服务是否正常运行: `curl http://localhost:3000/health`
- 检查防火墙和端口映射
- 检查 `GITLAB_WEBHOOK_SECRET` 配置是否正确
- 查看应用日志: `docker-compose logs -f webhook-server`

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart webhook-server

# 完全重建并启动
docker-compose up -d --build
```

### 清理数据

```bash
# 停止并删除容器
docker-compose down

# 删除容器和卷（会删除数据库数据）
docker-compose down -v

# 删除镜像
docker-compose down --rmi all
```

## 性能优化

### 资源限制

在生产环境中，建议设置资源限制（已在 `docker-compose.prod.yml` 中配置）：

```yaml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 1G
```

### 数据库优化

MySQL 配置优化（已在 `docker-compose.prod.yml` 中配置）：

```yaml
command:
  - --max_connections=200
  - --innodb_buffer_pool_size=512M
```

## 安全建议

1. **使用强密码**: 确保数据库密码足够复杂
2. **限制网络访问**: 生产环境不要暴露数据库端口
3. **使用 HTTPS**: 通过反向代理配置 SSL/TLS
4. **定期备份**: 配置数据库定期备份
5. **更新镜像**: 定期更新 Docker 镜像到最新版本

## 监控

### 日志监控

使用 Docker 日志驱动或日志收集工具（如 ELK、Loki）收集日志。

### 指标监控

可以集成 Prometheus 等监控工具收集应用指标。

---

**提示**: 更多部署相关问题，请参考项目 README 或提交 Issue。

