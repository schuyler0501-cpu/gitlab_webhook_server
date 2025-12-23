# 数据库配置指南

## 📊 支持的数据库

项目支持以下数据库：

- **MySQL** (默认) - 推荐用于生产环境
- **PostgreSQL** - 支持完整功能

## 🔧 配置方式

### 环境变量配置

在 `.env` 文件中配置数据库连接信息：

#### MySQL 配置（默认）

```env
# 数据库类型
DB_TYPE=mysql

# MySQL 连接配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gitlab_webhook
DB_CHARSET=utf8mb4
DB_TIMEZONE=Asia/Shanghai
```

#### PostgreSQL 配置

```env
# 数据库类型
DB_TYPE=postgresql

# PostgreSQL 连接配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gitlab_webhook
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Shanghai
```

## 🚀 快速开始

### 1. MySQL 设置

#### 创建数据库

```sql
CREATE DATABASE gitlab_webhook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 执行迁移

```bash
mysql -u root -p gitlab_webhook < migrations/001_create_tables_mysql.sql
```

或使用 MySQL 客户端：

```bash
mysql -u root -p
USE gitlab_webhook;
SOURCE migrations/001_create_tables_mysql.sql;
```

### 2. PostgreSQL 设置

#### 创建数据库

```sql
CREATE DATABASE gitlab_webhook;
```

#### 执行迁移

```bash
psql -U postgres -d gitlab_webhook -f migrations/001_create_tables.sql
```

### 3. 使用 GORM 自动迁移

项目支持 GORM 自动迁移，首次启动时会自动创建表结构。只需确保：

1. 数据库已创建
2. 配置正确的连接信息
3. 启动应用

## 📝 配置说明

### 数据库类型 (DB_TYPE)

- `mysql` - 使用 MySQL（默认）
- `postgresql` 或 `postgres` - 使用 PostgreSQL

### MySQL 配置项

- `DB_HOST` - 数据库主机（默认: localhost）
- `DB_PORT` - 数据库端口（默认: 3306）
- `DB_USER` - 数据库用户名（默认: root）
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名称（默认: gitlab_webhook）
- `DB_CHARSET` - 字符集（默认: utf8mb4）
- `DB_TIMEZONE` - 时区（默认: Asia/Shanghai）

### PostgreSQL 配置项

- `DB_HOST` - 数据库主机（默认: localhost）
- `DB_PORT` - 数据库端口（默认: 5432）
- `DB_USER` - 数据库用户名（默认: postgres）
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名称（默认: gitlab_webhook）
- `DB_SSLMODE` - SSL 模式（默认: disable）
- `DB_TIMEZONE` - 时区（默认: Asia/Shanghai）

## 🔍 验证连接

启动应用后，检查日志输出：

```
[INFO] 数据库连接成功 type=mysql host=localhost database=gitlab_webhook
```

如果连接失败，检查：

1. 数据库服务是否运行
2. 连接信息是否正确
3. 防火墙设置
4. 用户权限

## 🛠️ 切换数据库

### 从 PostgreSQL 切换到 MySQL

1. 修改 `.env` 文件：
   ```env
   DB_TYPE=mysql
   DB_PORT=3306
   DB_USER=root
   ```

2. 执行 MySQL 迁移文件

3. 重启应用

### 从 MySQL 切换到 PostgreSQL

1. 修改 `.env` 文件：
   ```env
   DB_TYPE=postgresql
   DB_PORT=5432
   DB_USER=postgres
   ```

2. 执行 PostgreSQL 迁移文件

3. 重启应用

## ⚠️ 注意事项

1. **数据迁移**: 切换数据库类型时，需要手动迁移数据
2. **字符集**: MySQL 建议使用 `utf8mb4` 以支持完整的 Unicode
3. **时区**: 确保数据库和应用使用相同的时区设置
4. **连接池**: 应用自动配置连接池（最大 100 个连接）
5. **自动迁移**: GORM 会自动创建/更新表结构，但不会删除字段

## 📚 相关文档

- [数据库设计文档](./database_design.md) - 详细的表结构设计
- [迁移文件](../migrations/) - SQL 迁移脚本

