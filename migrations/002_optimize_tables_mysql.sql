-- MySQL 数据库优化迁移文件：添加缺失字段和聚合统计表
-- 创建时间: 2024-12-24

-- 1. 为 commits 表添加缺失字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS project_id INT,
ADD COLUMN IF NOT EXISTS branch VARCHAR(255),
ADD COLUMN IF NOT EXISTS title VARCHAR(255),
ADD COLUMN IF NOT EXISTS committer_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS committer_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS authored_date DATETIME,
ADD COLUMN IF NOT EXISTS committed_date DATETIME;

-- 2. 更新现有记录的 committed_date（如果为空，使用 timestamp）
UPDATE commits 
SET committed_date = timestamp 
WHERE committed_date IS NULL;

UPDATE commits 
SET authored_date = timestamp 
WHERE authored_date IS NULL;

-- 3. 优化唯一索引：支持同一 commit 在不同项目中
-- 先删除旧的唯一索引
ALTER TABLE commits DROP INDEX IF EXISTS commits_commit_id_key;
ALTER TABLE commits DROP INDEX IF EXISTS idx_commits_commit_id;

-- 创建新的唯一索引（commit_id + project_id）
-- 注意：MySQL 不支持 COALESCE 在索引中，使用 IFNULL
CREATE UNIQUE INDEX IF NOT EXISTS idx_commits_commit_project 
ON commits(commit_id, IFNULL(project_id, 0));

-- 添加 project_id 索引
CREATE INDEX IF NOT EXISTS idx_commits_project_id ON commits(project_id);
CREATE INDEX IF NOT EXISTS idx_commits_branch ON commits(branch);
CREATE INDEX IF NOT EXISTS idx_commits_committed_date ON commits(committed_date);

-- 4. 创建成员贡献聚合统计表（提高查询性能）
CREATE TABLE IF NOT EXISTS member_contributions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_email VARCHAR(255) NOT NULL,
    member_name VARCHAR(255),
    project_id INT,
    project_name VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    commit_count INT DEFAULT 0,
    additions INT DEFAULT 0,
    deletions INT DEFAULT 0,
    net_lines INT AS (additions - deletions) STORED,
    total_changes INT AS (additions + deletions) STORED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_member_contributions (member_email, project_id, start_date, end_date),
    INDEX idx_member_contributions_email (member_email),
    INDEX idx_member_contributions_project (project_id),
    INDEX idx_member_contributions_period (start_date, end_date),
    INDEX idx_member_contributions_email_period (member_email, start_date, end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. 创建成员语言统计表（提高语言统计查询性能）
CREATE TABLE IF NOT EXISTS member_language_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_email VARCHAR(255) NOT NULL,
    language VARCHAR(100) NOT NULL,
    lines_added INT DEFAULT 0,
    lines_removed INT DEFAULT 0,
    file_count INT DEFAULT 0,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    project_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_member_language_stats (member_email, language, period_start, period_end, project_id),
    INDEX idx_member_language_stats_email (member_email),
    INDEX idx_member_language_stats_language (language),
    INDEX idx_member_language_stats_period (period_start, period_end),
    INDEX idx_member_language_stats_email_period (member_email, period_start, period_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

