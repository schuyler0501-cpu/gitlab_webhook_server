-- MySQL 数据库迁移文件：创建提交记录相关表
-- 创建时间: 2024-01-15

-- 1. 创建 commits 表 - 提交记录主表
CREATE TABLE IF NOT EXISTS commits (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    commit_id VARCHAR(255) NOT NULL UNIQUE,
    message TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    author VARCHAR(255) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    url TEXT,
    project_name VARCHAR(255) NOT NULL,
    project_path VARCHAR(500) NOT NULL,
    total_added_lines INT DEFAULT 0,
    total_removed_lines INT DEFAULT 0,
    total_changed_files INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_commits_author_email (author_email),
    INDEX idx_commits_timestamp (timestamp),
    INDEX idx_commits_project_name (project_name),
    INDEX idx_commits_commit_id (commit_id),
    INDEX idx_commits_author_timestamp (author_email, timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. 创建 commit_files 表 - 文件变更详情表
CREATE TABLE IF NOT EXISTS commit_files (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    commit_id BIGINT UNSIGNED NOT NULL,
    file_path TEXT NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_extension VARCHAR(50),
    change_type VARCHAR(20) NOT NULL,
    added_lines INT DEFAULT 0,
    removed_lines INT DEFAULT 0,
    language VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_commit_files_commit_id (commit_id),
    INDEX idx_commit_files_language (language),
    INDEX idx_commit_files_change_type (change_type),
    CONSTRAINT chk_change_type CHECK (change_type IN ('added', 'modified', 'removed')),
    CONSTRAINT fk_commit_files_commit FOREIGN KEY (commit_id) REFERENCES commits(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. 创建 commit_languages 表 - 语言统计表
CREATE TABLE IF NOT EXISTS commit_languages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    commit_id BIGINT UNSIGNED NOT NULL,
    language VARCHAR(50) NOT NULL,
    added_lines INT DEFAULT 0,
    removed_lines INT DEFAULT 0,
    file_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_commit_language (commit_id, language),
    INDEX idx_commit_languages_commit_id (commit_id),
    INDEX idx_commit_languages_language (language),
    CONSTRAINT fk_commit_languages_commit FOREIGN KEY (commit_id) REFERENCES commits(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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

