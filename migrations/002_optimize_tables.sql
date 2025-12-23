-- 数据库优化迁移文件：添加缺失字段和聚合统计表
-- 创建时间: 2024-12-24
-- 注意: 这是 PostgreSQL 版本，MySQL 版本请使用 002_optimize_tables_mysql.sql

-- 1. 为 commits 表添加缺失字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS project_id INTEGER,
ADD COLUMN IF NOT EXISTS branch VARCHAR(255),
ADD COLUMN IF NOT EXISTS title VARCHAR(255),
ADD COLUMN IF NOT EXISTS committer_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS committer_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS authored_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS committed_date TIMESTAMP;

-- 2. 更新现有记录的 committed_date（如果为空，使用 timestamp）
UPDATE commits 
SET committed_date = timestamp 
WHERE committed_date IS NULL;

UPDATE commits 
SET authored_date = timestamp 
WHERE authored_date IS NULL;

-- 3. 优化唯一索引：支持同一 commit 在不同项目中
-- 先删除旧的唯一索引
DROP INDEX IF EXISTS commits_commit_id_key;
DROP INDEX IF EXISTS idx_commits_commit_id;

-- 创建新的唯一索引（commit_id + project_id）
-- 注意：如果 project_id 可能为空，需要先处理
CREATE UNIQUE INDEX IF NOT EXISTS idx_commits_commit_project 
ON commits(commit_id, COALESCE(project_id, 0));

-- 添加 project_id 索引
CREATE INDEX IF NOT EXISTS idx_commits_project_id ON commits(project_id);
CREATE INDEX IF NOT EXISTS idx_commits_branch ON commits(branch);
CREATE INDEX IF NOT EXISTS idx_commits_committed_date ON commits(committed_date);

-- 4. 创建成员贡献聚合统计表（提高查询性能）
CREATE TABLE IF NOT EXISTS member_contributions (
    id BIGSERIAL PRIMARY KEY,
    member_email VARCHAR(255) NOT NULL,
    member_name VARCHAR(255),
    project_id INTEGER,
    project_name VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    commit_count INTEGER DEFAULT 0,
    additions INTEGER DEFAULT 0,
    deletions INTEGER DEFAULT 0,
    net_lines INTEGER GENERATED ALWAYS AS (additions - deletions) STORED,
    total_changes INTEGER GENERATED ALWAYS AS (additions + deletions) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(member_email, project_id, start_date, end_date)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_member_contributions_email ON member_contributions(member_email);
CREATE INDEX IF NOT EXISTS idx_member_contributions_project ON member_contributions(project_id);
CREATE INDEX IF NOT EXISTS idx_member_contributions_period ON member_contributions(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_member_contributions_email_period ON member_contributions(member_email, start_date, end_date);

-- 5. 创建成员语言统计表（提高语言统计查询性能）
CREATE TABLE IF NOT EXISTS member_language_stats (
    id BIGSERIAL PRIMARY KEY,
    member_email VARCHAR(255) NOT NULL,
    language VARCHAR(100) NOT NULL,
    lines_added INTEGER DEFAULT 0,
    lines_removed INTEGER DEFAULT 0,
    file_count INTEGER DEFAULT 0,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    project_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(member_email, language, period_start, period_end, project_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_member_language_stats_email ON member_language_stats(member_email);
CREATE INDEX IF NOT EXISTS idx_member_language_stats_language ON member_language_stats(language);
CREATE INDEX IF NOT EXISTS idx_member_language_stats_period ON member_language_stats(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_member_language_stats_email_period ON member_language_stats(member_email, period_start, period_end);

-- 添加注释
COMMENT ON TABLE member_contributions IS '成员贡献聚合统计表，用于提高统计查询性能';
COMMENT ON TABLE member_language_stats IS '成员语言统计表，用于提高语言统计查询性能';
COMMENT ON COLUMN commits.project_id IS 'GitLab 项目 ID';
COMMENT ON COLUMN commits.branch IS '提交所在分支';
COMMENT ON COLUMN commits.title IS '提交标题（message 第一行）';
COMMENT ON COLUMN commits.committer_name IS '提交者姓名（可能与作者不同）';
COMMENT ON COLUMN commits.committer_email IS '提交者邮箱';
COMMENT ON COLUMN commits.authored_date IS '代码编写时间';
COMMENT ON COLUMN commits.committed_date IS '代码提交时间';

