# 数据库设计文档

## 📊 数据库表结构

### 1. commits 表 - 提交记录主表

存储每次代码提交的基本信息。

```sql
CREATE TABLE commits (
    id BIGSERIAL PRIMARY KEY,
    commit_id VARCHAR(255) NOT NULL UNIQUE,
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    author VARCHAR(255) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    url TEXT,
    project_name VARCHAR(255) NOT NULL,
    project_path VARCHAR(500) NOT NULL,
    total_added_lines INTEGER DEFAULT 0,
    total_removed_lines INTEGER DEFAULT 0,
    total_changed_files INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_commits_author_email ON commits(author_email);
CREATE INDEX idx_commits_timestamp ON commits(timestamp);
CREATE INDEX idx_commits_project_name ON commits(project_name);
CREATE INDEX idx_commits_commit_id ON commits(commit_id);
```

**字段说明**:

- `id`: 主键，自增
- `commit_id`: GitLab 提交 ID，与 `project_id` 组成唯一索引
- `project_id`: GitLab 项目 ID
- `message`: 提交信息
- `title`: 提交标题（message 第一行）
- `timestamp`: 提交时间（保持向后兼容）
- `author`: 作者姓名
- `author_email`: 作者邮箱（用于统计）
- `committer_name/committer_email`: 提交者信息（可能与作者不同）
- `authored_date/committed_date`: 编写时间和提交时间
- `branch`: 提交所在分支
- `ref_protected`: 分支是否受保护
- `url`: 提交链接
- `project_name`: 项目名称
- `project_path`: 项目路径
- `project_description`: 项目描述
- `project_web_url`: 项目 Web URL
- `project_namespace`: 项目命名空间
- `project_visibility_level`: 项目可见性级别（0=private, 10=internal, 20=public）
- `project_default_branch`: 项目默认分支
- `project_git_ssh_url/project_git_http_url`: 项目 Git URL
- `repository_name`: 仓库名称
- `repository_url`: 仓库 URL
- `repository_description`: 仓库描述
- `repository_homepage`: 仓库主页
- `repository_git_ssh_url/repository_git_http_url`: 仓库 Git URL
- `repository_visibility_level`: 仓库可见性级别
- `before_sha/after_sha/checkout_sha`: 推送相关的 SHA
- `push_message`: 推送消息
- `total_commits_count`: 本次推送的总提交数
- `push_user_id/push_user_name/push_user_username/push_user_email`: 推送用户信息（推送者，可能与提交作者不同）
- `total_added_lines`: 总新增行数
- `total_removed_lines`: 总删除行数
- `total_changed_files`: 总变更文件数
- `created_at`: 记录创建时间
- `updated_at`: 记录更新时间

### 2. commit_files 表 - 文件变更详情表

存储每次提交中每个文件的详细变更信息。

```sql
CREATE TABLE commit_files (
    id BIGSERIAL PRIMARY KEY,
    commit_id BIGINT NOT NULL REFERENCES commits(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_extension VARCHAR(50),
    change_type VARCHAR(20) NOT NULL, -- 'added', 'modified', 'removed'
    added_lines INTEGER DEFAULT 0,
    removed_lines INTEGER DEFAULT 0,
    language VARCHAR(50), -- 编程语言，如 'go', 'java', 'python'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_commit_files_commit_id ON commit_files(commit_id);
CREATE INDEX idx_commit_files_language ON commit_files(language);
CREATE INDEX idx_commit_files_change_type ON commit_files(change_type);
```

**字段说明**:

- `id`: 主键，自增
- `commit_id`: 关联的提交 ID（外键）
- `file_path`: 文件完整路径
- `file_name`: 文件名
- `file_extension`: 文件扩展名
- `change_type`: 变更类型（added/modified/removed）
- `added_lines`: 新增行数
- `removed_lines`: 删除行数
- `language`: 编程语言
- `created_at`: 记录创建时间

### 3. commit_languages 表 - 语言统计表

存储每次提交中每种编程语言的代码行数统计。

```sql
CREATE TABLE commit_languages (
    id BIGSERIAL PRIMARY KEY,
    commit_id BIGINT NOT NULL REFERENCES commits(id) ON DELETE CASCADE,
    language VARCHAR(50) NOT NULL,
    added_lines INTEGER DEFAULT 0,
    removed_lines INTEGER DEFAULT 0,
    file_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(commit_id, language)
);

CREATE INDEX idx_commit_languages_commit_id ON commit_languages(commit_id);
CREATE INDEX idx_commit_languages_language ON commit_languages(language);
```

**字段说明**:

- `id`: 主键，自增
- `commit_id`: 关联的提交 ID（外键）
- `language`: 编程语言名称
- `added_lines`: 该语言新增行数
- `removed_lines`: 该语言删除行数
- `file_count`: 该语言文件数量
- `created_at`: 记录创建时间

## 📈 统计查询视图

### 成员提交统计视图

```sql
CREATE VIEW member_commit_stats AS
SELECT 
    author_email,
    DATE(timestamp) as commit_date,
    COUNT(*) as commit_count,
    SUM(total_added_lines) as total_added,
    SUM(total_removed_lines) as total_removed,
    SUM(total_changed_files) as total_files
FROM commits
GROUP BY author_email, DATE(timestamp);
```

### 语言统计视图

```sql
CREATE VIEW language_stats AS
SELECT 
    cl.language,
    COUNT(DISTINCT cl.commit_id) as commit_count,
    SUM(cl.added_lines) as total_added,
    SUM(cl.removed_lines) as total_removed,
    SUM(cl.file_count) as total_files
FROM commit_languages cl
GROUP BY cl.language;
```

## 🔍 常用查询示例

### 1. 查询成员在时间范围内的提交统计

```sql
SELECT 
    author_email,
    COUNT(*) as commit_count,
    SUM(total_added_lines) as total_added,
    SUM(total_removed_lines) as total_removed,
    SUM(total_changed_files) as total_files
FROM commits
WHERE author_email = 'user@example.com'
  AND timestamp >= '2024-01-01'
  AND timestamp < '2024-02-01'
GROUP BY author_email;
```

### 2. 查询成员的语言使用统计

```sql
SELECT 
    cl.language,
    SUM(cl.added_lines) as total_added,
    SUM(cl.removed_lines) as total_removed,
    SUM(cl.file_count) as total_files
FROM commit_languages cl
JOIN commits c ON cl.commit_id = c.id
WHERE c.author_email = 'user@example.com'
  AND c.timestamp >= '2024-01-01'
  AND c.timestamp < '2024-02-01'
GROUP BY cl.language
ORDER BY total_added DESC;
```

### 3. 查询项目的提交统计

```sql
SELECT 
    project_name,
    COUNT(*) as commit_count,
    COUNT(DISTINCT author_email) as member_count,
    SUM(total_added_lines) as total_added,
    SUM(total_removed_lines) as total_removed
FROM commits
WHERE project_name = 'project-name'
  AND timestamp >= '2024-01-01'
GROUP BY project_name;
```

## 🎯 设计原则

1. **规范化**: 使用三张表分离提交、文件变更和语言统计
2. **索引优化**: 为常用查询字段创建索引
3. **外键约束**: 使用外键确保数据完整性
4. **级联删除**: 删除提交时自动删除关联记录
5. **时间戳**: 记录创建和更新时间，便于追踪

## 📝 注意事项

1. **GitLab Webhook 数据**: 需要从 Webhook payload 中提取代码行数信息
2. **语言识别**: 根据文件扩展名识别编程语言
3. **性能优化**: 对于大量数据，考虑分区表
4. **数据一致性**: 使用事务确保数据一致性
