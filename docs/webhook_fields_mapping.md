# GitLab Webhook 字段映射文档

## 📋 字段映射表

本文档说明 GitLab Webhook JSON 结构与数据库表字段的完整映射关系。

### 顶层字段（Push Event）

| Webhook JSON 字段 | 数据库字段 | 说明 | 是否必需 |
|------------------|-----------|------|---------|
| `object_kind` | - | 事件类型（"push"） | 否（用于路由） |
| `event_name` | - | 事件名称（"push"） | 否（用于路由） |
| `before` | `before_sha` | 推送前的 commit SHA | 是 |
| `after` | `after_sha` | 推送后的 commit SHA | 是 |
| `ref` | `branch` | 分支引用（解析为分支名） | 是 |
| `ref_protected` | `ref_protected` | 分支是否受保护 | 是 |
| `checkout_sha` | `checkout_sha` | checkout SHA | 是 |
| `message` | `push_message` | 推送消息 | 否 |
| `user_id` | `push_user_id` | 推送用户 ID | 是 |
| `user_name` | `push_user_name` | 推送用户名称 | 是 |
| `user_username` | `push_user_username` | 推送用户用户名 | 是 |
| `user_email` | `push_user_email` | 推送用户邮箱 | 是 |
| `project_id` | `project_id` | 项目 ID | 是 |
| `total_commits_count` | `total_commits_count` | 本次推送的总提交数 | 是 |

### project 对象字段

| Webhook JSON 字段 | 数据库字段 | 说明 | 是否必需 |
|------------------|-----------|------|---------|
| `project.id` | `project_id` | 项目 ID | 是 |
| `project.name` | `project_name` | 项目名称 | 是 |
| `project.path_with_namespace` | `project_path` | 项目路径（含命名空间） | 是 |
| `project.description` | `project_description` | 项目描述 | 否 |
| `project.web_url` | `project_web_url` | 项目 Web URL | 否 |
| `project.namespace` | `project_namespace` | 项目命名空间 | 是 |
| `project.visibility_level` | `project_visibility_level` | 项目可见性级别 | 是 |
| `project.default_branch` | `project_default_branch` | 项目默认分支 | 否 |
| `project.git_ssh_url` | `project_git_ssh_url` | 项目 Git SSH URL | 否 |
| `project.git_http_url` | `project_git_http_url` | 项目 Git HTTP URL | 否 |

### repository 对象字段

| Webhook JSON 字段 | 数据库字段 | 说明 | 是否必需 |
|------------------|-----------|------|---------|
| `repository.name` | `repository_name` | 仓库名称 | 否 |
| `repository.url` | `repository_url` | 仓库 URL | 否 |
| `repository.description` | `repository_description` | 仓库描述 | 否 |
| `repository.homepage` | `repository_homepage` | 仓库主页 | 否 |
| `repository.git_ssh_url` | `repository_git_ssh_url` | 仓库 Git SSH URL | 否 |
| `repository.git_http_url` | `repository_git_http_url` | 仓库 Git HTTP URL | 否 |
| `repository.visibility_level` | `repository_visibility_level` | 仓库可见性级别 | 否 |

### commits 数组中的字段（每个提交）

| Webhook JSON 字段 | 数据库字段 | 说明 | 是否必需 |
|------------------|-----------|------|---------|
| `commits[].id` | `commit_id` | 提交 ID（SHA） | 是 |
| `commits[].message` | `message` | 提交信息 | 是 |
| `commits[].title` | `title` | 提交标题（message 第一行） | 是 |
| `commits[].timestamp` | `timestamp`, `committed_date` | 提交时间 | 是 |
| `commits[].url` | `url` | 提交链接 | 是 |
| `commits[].author.name` | `author` | 作者姓名 | 是 |
| `commits[].author.email` | `author_email` | 作者邮箱 | 是 |
| `commits[].committer.name` | `committer_name` | 提交者姓名 | 否（默认同作者） |
| `commits[].committer.email` | `committer_email` | 提交者邮箱 | 否（默认同作者） |
| `commits[].added` | `commit_files` (change_type='added') | 新增文件列表 | 是 |
| `commits[].modified` | `commit_files` (change_type='modified') | 修改文件列表 | 是 |
| `commits[].removed` | `commit_files` (change_type='removed') | 删除文件列表 | 是 |

## 🔍 字段分类

### 核心字段（必需）
- 提交标识：`commit_id`, `project_id`
- 提交信息：`message`, `title`, `timestamp`
- 作者信息：`author`, `author_email`
- 项目信息：`project_name`, `project_path`
- 分支信息：`branch`

### 扩展字段（推荐）
- 推送用户信息：`push_user_id`, `push_user_name`, `push_user_username`
- 分支保护：`ref_protected`
- 项目扩展：`project_namespace`, `project_visibility_level`
- 推送 SHA：`before_sha`, `after_sha`, `checkout_sha`

### 可选字段（按需）
- 项目描述：`project_description`
- 项目 URL：`project_web_url`, `project_git_ssh_url`, `project_git_http_url`
- 仓库信息：`repository_name`, `repository_url`, `repository_description`
- 推送消息：`push_message`

## 📊 数据流向

```
GitLab Webhook JSON
    ↓
parsePushInfo() - 解析推送级别信息
    ↓
parseCommit() - 解析每个提交
    ↓
CommitRecord - 内存模型
    ↓
RecordCommit() - 保存到数据库
    ↓
Commit (数据库模型)
```

## 🎯 使用场景

### 1. 效能度量
- **推送用户信息**：区分推送者和提交作者（可能不同）
- **分支保护状态**：识别重要分支的提交
- **可见性级别**：区分公开/内部/私有项目

### 2. 统计分析
- **命名空间**：按组织/团队统计
- **项目描述**：项目分类和标签
- **仓库信息**：仓库级别的统计

### 3. 审计追踪
- **推送 SHA**：完整的推送链路追踪
- **推送消息**：推送操作的说明
- **总提交数**：批量推送的统计

## ⚠️ 注意事项

1. **推送用户 vs 提交作者**：
   - 推送用户（`push_user_*`）：执行推送操作的用户
   - 提交作者（`author_*`）：代码的实际作者
   - 两者可能不同（如合并操作）

2. **分支信息**：
   - `ref` 格式：`refs/heads/master`
   - 解析后存储为：`master`

3. **可见性级别**：
   - `0` = Private（私有）
   - `10` = Internal（内部）
   - `20` = Public（公开）

4. **时间字段**：
   - `timestamp`：保持向后兼容
   - `committed_date`：实际提交时间
   - `authored_date`：代码编写时间（通常相同）

5. **唯一性**：
   - 使用 `(commit_id, project_id)` 作为唯一索引
   - 支持同一 commit 在不同项目中

## 📝 迁移说明

执行迁移文件以添加新字段：

```bash
# PostgreSQL
psql -U username -d database_name -f migrations/003_add_webhook_fields.sql

# MySQL
mysql -u username -p database_name < migrations/003_add_webhook_fields_mysql.sql
```

**注意**：MySQL 迁移文件不包含 `IF NOT EXISTS`，如果字段已存在会报错。建议先检查或使用 GORM 的 AutoMigrate。

