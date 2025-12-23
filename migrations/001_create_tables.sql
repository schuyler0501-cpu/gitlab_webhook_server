-- PostgreSQL 数据库迁移文件：创建提交记录相关表
-- 创建时间: 2024-01-15
-- 注意: 这是 PostgreSQL 版本，MySQL 版本请使用 001_create_tables_mysql.sql

-- 1. 创建 commits 表 - 提交记录主表
CREATE TABLE IF NOT EXISTS commits (
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

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_commits_author_email ON commits(author_email);
CREATE INDEX IF NOT EXISTS idx_commits_timestamp ON commits(timestamp);
CREATE INDEX IF NOT EXISTS idx_commits_project_name ON commits(project_name);
CREATE INDEX IF NOT EXISTS idx_commits_commit_id ON commits(commit_id);
CREATE INDEX IF NOT EXISTS idx_commits_author_timestamp ON commits(author_email, timestamp);

-- 2. 创建 commit_files 表 - 文件变更详情表
CREATE TABLE IF NOT EXISTS commit_files (
    id BIGSERIAL PRIMARY KEY,
    commit_id BIGINT NOT NULL REFERENCES commits(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_extension VARCHAR(50),
    change_type VARCHAR(20) NOT NULL CHECK (change_type IN ('added', 'modified', 'removed')),
    added_lines INTEGER DEFAULT 0,
    removed_lines INTEGER DEFAULT 0,
    language VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_commit_files_commit_id ON commit_files(commit_id);
CREATE INDEX IF NOT EXISTS idx_commit_files_language ON commit_files(language);
CREATE INDEX IF NOT EXISTS idx_commit_files_change_type ON commit_files(change_type);

-- 3. 创建 commit_languages 表 - 语言统计表
CREATE TABLE IF NOT EXISTS commit_languages (
    id BIGSERIAL PRIMARY KEY,
    commit_id BIGINT NOT NULL REFERENCES commits(id) ON DELETE CASCADE,
    language VARCHAR(50) NOT NULL,
    added_lines INTEGER DEFAULT 0,
    removed_lines INTEGER DEFAULT 0,
    file_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(commit_id, language)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_commit_languages_commit_id ON commit_languages(commit_id);
CREATE INDEX IF NOT EXISTS idx_commit_languages_language ON commit_languages(language);

-- 4. 创建统计视图

-- 成员提交统计视图
CREATE OR REPLACE VIEW member_commit_stats AS
SELECT 
    author_email,
    DATE(timestamp) as commit_date,
    COUNT(*) as commit_count,
    SUM(total_added_lines) as total_added,
    SUM(total_removed_lines) as total_removed,
    SUM(total_changed_files) as total_files
FROM commits
GROUP BY author_email, DATE(timestamp);

-- 语言统计视图
CREATE OR REPLACE VIEW language_stats AS
SELECT 
    cl.language,
    COUNT(DISTINCT cl.commit_id) as commit_count,
    SUM(cl.added_lines) as total_added,
    SUM(cl.removed_lines) as total_removed,
    SUM(cl.file_count) as total_files
FROM commit_languages cl
GROUP BY cl.language;

-- 添加注释
COMMENT ON TABLE commits IS '代码提交记录主表';
COMMENT ON TABLE commit_files IS '文件变更详情表';
COMMENT ON TABLE commit_languages IS '语言统计表';
COMMENT ON COLUMN commits.commit_id IS 'GitLab 提交 ID，唯一标识';
COMMENT ON COLUMN commits.total_added_lines IS '本次提交总新增行数';
COMMENT ON COLUMN commits.total_removed_lines IS '本次提交总删除行数';
COMMENT ON COLUMN commit_files.change_type IS '变更类型：added/modified/removed';
COMMENT ON COLUMN commit_files.language IS '编程语言，根据文件扩展名识别';

