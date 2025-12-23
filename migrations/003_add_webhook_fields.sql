-- 数据库迁移文件：添加 GitLab Webhook 完整字段支持
-- 创建时间: 2024-12-24
-- 注意: 这是 PostgreSQL 版本，MySQL 版本请使用 003_add_webhook_fields_mysql.sql

-- 1. 为 commits 表添加推送用户信息字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS push_user_id INTEGER,
ADD COLUMN IF NOT EXISTS push_user_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS push_user_username VARCHAR(255),
ADD COLUMN IF NOT EXISTS push_user_email VARCHAR(255);

-- 2. 添加分支保护状态字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS ref_protected BOOLEAN DEFAULT FALSE;

-- 3. 添加推送相关字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS before_sha VARCHAR(40),
ADD COLUMN IF NOT EXISTS after_sha VARCHAR(40),
ADD COLUMN IF NOT EXISTS checkout_sha VARCHAR(40),
ADD COLUMN IF NOT EXISTS push_message TEXT,
ADD COLUMN IF NOT EXISTS total_commits_count INTEGER DEFAULT 0;

-- 4. 添加项目扩展信息字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS project_description TEXT,
ADD COLUMN IF NOT EXISTS project_web_url TEXT,
ADD COLUMN IF NOT EXISTS project_namespace VARCHAR(255),
ADD COLUMN IF NOT EXISTS project_visibility_level INTEGER,
ADD COLUMN IF NOT EXISTS project_default_branch VARCHAR(255),
ADD COLUMN IF NOT EXISTS project_git_ssh_url TEXT,
ADD COLUMN IF NOT EXISTS project_git_http_url TEXT;

-- 5. 添加仓库信息字段
ALTER TABLE commits 
ADD COLUMN IF NOT EXISTS repository_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS repository_url TEXT,
ADD COLUMN IF NOT EXISTS repository_description TEXT,
ADD COLUMN IF NOT EXISTS repository_homepage TEXT,
ADD COLUMN IF NOT EXISTS repository_git_ssh_url TEXT,
ADD COLUMN IF NOT EXISTS repository_git_http_url TEXT,
ADD COLUMN IF NOT EXISTS repository_visibility_level INTEGER;

-- 6. 创建索引
CREATE INDEX IF NOT EXISTS idx_commits_push_user_id ON commits(push_user_id);
CREATE INDEX IF NOT EXISTS idx_commits_push_user_username ON commits(push_user_username);
CREATE INDEX IF NOT EXISTS idx_commits_ref_protected ON commits(ref_protected);
CREATE INDEX IF NOT EXISTS idx_commits_project_namespace ON commits(project_namespace);
CREATE INDEX IF NOT EXISTS idx_commits_project_visibility ON commits(project_visibility_level);

-- 添加注释
COMMENT ON COLUMN commits.push_user_id IS '推送用户 ID（GitLab user_id）';
COMMENT ON COLUMN commits.push_user_name IS '推送用户名称';
COMMENT ON COLUMN commits.push_user_username IS '推送用户用户名';
COMMENT ON COLUMN commits.push_user_email IS '推送用户邮箱';
COMMENT ON COLUMN commits.ref_protected IS '分支是否受保护';
COMMENT ON COLUMN commits.before_sha IS '推送前的 commit SHA';
COMMENT ON COLUMN commits.after_sha IS '推送后的 commit SHA';
COMMENT ON COLUMN commits.checkout_sha IS 'checkout SHA';
COMMENT ON COLUMN commits.push_message IS '推送消息';
COMMENT ON COLUMN commits.total_commits_count IS '本次推送的总提交数';
COMMENT ON COLUMN commits.project_namespace IS '项目命名空间';
COMMENT ON COLUMN commits.project_visibility_level IS '项目可见性级别（0=private, 10=internal, 20=public）';

