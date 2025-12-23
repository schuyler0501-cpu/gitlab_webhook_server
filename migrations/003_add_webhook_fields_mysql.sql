-- MySQL 数据库迁移文件：添加 GitLab Webhook 完整字段支持
-- 创建时间: 2024-12-24

-- 1. 为 commits 表添加推送用户信息字段
-- MySQL 不支持 IF NOT EXISTS，需要手动检查或使用存储过程
ALTER TABLE commits 
ADD COLUMN push_user_id INT,
ADD COLUMN push_user_name VARCHAR(255),
ADD COLUMN push_user_username VARCHAR(255),
ADD COLUMN push_user_email VARCHAR(255);

-- 2. 添加分支保护状态字段
ALTER TABLE commits 
ADD COLUMN ref_protected BOOLEAN DEFAULT FALSE;

-- 3. 添加推送相关字段
ALTER TABLE commits 
ADD COLUMN before_sha VARCHAR(40),
ADD COLUMN after_sha VARCHAR(40),
ADD COLUMN checkout_sha VARCHAR(40),
ADD COLUMN push_message TEXT,
ADD COLUMN total_commits_count INT DEFAULT 0;

-- 4. 添加项目扩展信息字段
ALTER TABLE commits 
ADD COLUMN project_description TEXT,
ADD COLUMN project_web_url TEXT,
ADD COLUMN project_namespace VARCHAR(255),
ADD COLUMN project_visibility_level INT,
ADD COLUMN project_default_branch VARCHAR(255),
ADD COLUMN project_git_ssh_url TEXT,
ADD COLUMN project_git_http_url TEXT;

-- 5. 添加仓库信息字段
ALTER TABLE commits 
ADD COLUMN repository_name VARCHAR(255),
ADD COLUMN repository_url TEXT,
ADD COLUMN repository_description TEXT,
ADD COLUMN repository_homepage TEXT,
ADD COLUMN repository_git_ssh_url TEXT,
ADD COLUMN repository_git_http_url TEXT,
ADD COLUMN repository_visibility_level INT;

-- 6. 创建索引
CREATE INDEX IF NOT EXISTS idx_commits_push_user_id ON commits(push_user_id);
CREATE INDEX IF NOT EXISTS idx_commits_push_user_username ON commits(push_user_username);
CREATE INDEX IF NOT EXISTS idx_commits_ref_protected ON commits(ref_protected);
CREATE INDEX IF NOT EXISTS idx_commits_project_namespace ON commits(project_namespace);
CREATE INDEX IF NOT EXISTS idx_commits_project_visibility ON commits(project_visibility_level);

