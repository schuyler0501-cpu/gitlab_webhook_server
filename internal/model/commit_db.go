package model

import (
	"time"

	"gorm.io/gorm"
)

// Commit 提交记录数据库模型
type Commit struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommitID         string    `gorm:"type:varchar(255);not null;index:idx_commits_commit_project,unique" json:"commit_id"`
	ProjectID        *int      `gorm:"type:integer;index;index:idx_commits_commit_project,unique" json:"project_id"`
	Message          string    `gorm:"type:text;not null" json:"message"`
	Title            string    `gorm:"type:varchar(255)" json:"title"`
	Timestamp        time.Time `gorm:"type:timestamp;not null;index" json:"timestamp"` // 保持向后兼容
	Author           string    `gorm:"type:varchar(255);not null" json:"author"`
	AuthorEmail      string    `gorm:"type:varchar(255);not null;index" json:"author_email"`
	CommitterName    string    `gorm:"type:varchar(255)" json:"committer_name"`
	CommitterEmail   string    `gorm:"type:varchar(255)" json:"committer_email"`
	AuthoredDate     *time.Time `gorm:"type:timestamp" json:"authored_date"`
	CommittedDate    *time.Time `gorm:"type:timestamp;index" json:"committed_date"`
	Branch           string    `gorm:"type:varchar(255);index" json:"branch"`
	RefProtected    *bool     `gorm:"type:boolean;default:false;index" json:"ref_protected"`
	URL              string    `gorm:"type:text" json:"url"`
	ProjectName      string    `gorm:"type:varchar(255);not null;index" json:"project_name"`
	ProjectPath      string    `gorm:"type:varchar(500);not null" json:"project_path"`
	// 推送用户信息（推送者，可能与提交作者不同）
	PushUserID       *int      `gorm:"type:integer;index" json:"push_user_id"`
	PushUserName     string    `gorm:"type:varchar(255)" json:"push_user_name"`
	PushUserUsername string    `gorm:"type:varchar(255);index" json:"push_user_username"`
	PushUserEmail    string    `gorm:"type:varchar(255)" json:"push_user_email"`
	// 推送相关字段
	BeforeSHA        string    `gorm:"type:varchar(40)" json:"before_sha"`
	AfterSHA         string    `gorm:"type:varchar(40)" json:"after_sha"`
	CheckoutSHA      string    `gorm:"type:varchar(40)" json:"checkout_sha"`
	PushMessage      string    `gorm:"type:text" json:"push_message"`
	TotalCommitsCount int      `gorm:"type:integer;default:0" json:"total_commits_count"`
	// 项目扩展信息
	ProjectDescription   string    `gorm:"type:text" json:"project_description"`
	ProjectWebURL         string    `gorm:"type:text" json:"project_web_url"`
	ProjectNamespace      string    `gorm:"type:varchar(255);index" json:"project_namespace"`
	ProjectVisibilityLevel *int     `gorm:"type:integer;index" json:"project_visibility_level"`
	ProjectDefaultBranch  string    `gorm:"type:varchar(255)" json:"project_default_branch"`
	ProjectGitSSHURL      string    `gorm:"type:text" json:"project_git_ssh_url"`
	ProjectGitHTTPURL     string    `gorm:"type:text" json:"project_git_http_url"`
	// 仓库信息
	RepositoryName         string    `gorm:"type:varchar(255)" json:"repository_name"`
	RepositoryURL          string    `gorm:"type:text" json:"repository_url"`
	RepositoryDescription  string    `gorm:"type:text" json:"repository_description"`
	RepositoryHomepage     string    `gorm:"type:text" json:"repository_homepage"`
	RepositoryGitSSHURL    string    `gorm:"type:text" json:"repository_git_ssh_url"`
	RepositoryGitHTTPURL   string    `gorm:"type:text" json:"repository_git_http_url"`
	RepositoryVisibilityLevel *int   `gorm:"type:integer" json:"repository_visibility_level"`
	TotalAddedLines  int       `gorm:"type:integer;default:0" json:"total_added_lines"`
	TotalRemovedLines int      `gorm:"type:integer;default:0" json:"total_removed_lines"`
	TotalChangedFiles int      `gorm:"type:integer;default:0" json:"total_changed_files"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Files    []CommitFile    `gorm:"foreignKey:CommitID;references:ID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
	Languages []CommitLanguage `gorm:"foreignKey:CommitID;references:ID;constraint:OnDelete:CASCADE" json:"languages,omitempty"`
}

// TableName 指定表名
func (Commit) TableName() string {
	return "commits"
}

// CommitFile 文件变更数据库模型
type CommitFile struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommitID     uint64    `gorm:"type:bigint;not null;index" json:"commit_id"`
	FilePath     string    `gorm:"type:text;not null" json:"file_path"`
	FileName     string    `gorm:"type:varchar(500);not null" json:"file_name"`
	FileExtension string   `gorm:"type:varchar(50)" json:"file_extension"`
	ChangeType   string    `gorm:"type:varchar(20);not null;index;check:change_type IN ('added','modified','removed')" json:"change_type"`
	AddedLines   int       `gorm:"type:integer;default:0" json:"added_lines"`
	RemovedLines int       `gorm:"type:integer;default:0" json:"removed_lines"`
	Language     string    `gorm:"type:varchar(50);index" json:"language"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (CommitFile) TableName() string {
	return "commit_files"
}

// CommitLanguage 语言统计数据库模型
type CommitLanguage struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommitID     uint64    `gorm:"type:bigint;not null;index" json:"commit_id"`
	Language     string    `gorm:"type:varchar(50);not null;index" json:"language"`
	AddedLines   int       `gorm:"type:integer;default:0" json:"added_lines"`
	RemovedLines int       `gorm:"type:integer;default:0" json:"removed_lines"`
	FileCount    int       `gorm:"type:integer;default:0" json:"file_count"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (CommitLanguage) TableName() string {
	return "commit_languages"
}

// BeforeCreate 创建前钩子
func (c *Commit) BeforeCreate(tx *gorm.DB) error {
	// 确保 (commit_id, project_id) 组合唯一
	var count int64
	query := tx.Model(&Commit{}).Where("commit_id = ?", c.CommitID)
	if c.ProjectID != nil {
		query = query.Where("project_id = ?", *c.ProjectID)
	} else {
		query = query.Where("project_id IS NULL")
	}
	query.Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

