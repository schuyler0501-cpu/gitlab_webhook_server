package model

import "time"

// CommitRecord 代码提交记录
type CommitRecord struct {
	CommitID       string            `json:"commit_id"`
	ProjectID      *int              `json:"project_id,omitempty"`
	Message        string            `json:"message"`
	Title          string            `json:"title,omitempty"`
	Timestamp      string            `json:"timestamp"`
	Author         string            `json:"author"`
	AuthorEmail    string            `json:"author_email"`
	CommitterName  string            `json:"committer_name,omitempty"`
	CommitterEmail string            `json:"committer_email,omitempty"`
	AuthoredDate   *time.Time        `json:"authored_date,omitempty"`
	CommittedDate  *time.Time        `json:"committed_date,omitempty"`
	Branch         string            `json:"branch,omitempty"`
	RefProtected   *bool             `json:"ref_protected,omitempty"`
	URL            string            `json:"url"`
	ProjectName    string            `json:"project_name"`
	ProjectPath    string            `json:"project_path"`
	// 项目扩展信息
	ProjectDescription    string    `json:"project_description,omitempty"`
	ProjectWebURL         string    `json:"project_web_url,omitempty"`
	ProjectNamespace      string    `json:"project_namespace,omitempty"`
	ProjectVisibilityLevel *int     `json:"project_visibility_level,omitempty"`
	ProjectDefaultBranch  string    `json:"project_default_branch,omitempty"`
	ProjectGitSSHURL      string    `json:"project_git_ssh_url,omitempty"`
	ProjectGitHTTPURL     string    `json:"project_git_http_url,omitempty"`
	// 仓库信息
	RepositoryName         string    `json:"repository_name,omitempty"`
	RepositoryURL          string    `json:"repository_url,omitempty"`
	RepositoryDescription  string    `json:"repository_description,omitempty"`
	RepositoryHomepage     string    `json:"repository_homepage,omitempty"`
	RepositoryGitSSHURL    string    `json:"repository_git_ssh_url,omitempty"`
	RepositoryGitHTTPURL   string    `json:"repository_git_http_url,omitempty"`
	RepositoryVisibilityLevel *int   `json:"repository_visibility_level,omitempty"`
	// 推送信息
	BeforeSHA              string    `json:"before_sha,omitempty"`
	AfterSHA               string    `json:"after_sha,omitempty"`
	CheckoutSHA            string    `json:"checkout_sha,omitempty"`
	PushMessage            string    `json:"push_message,omitempty"`
	TotalCommitsCount      int       `json:"total_commits_count,omitempty"`
	// 推送用户信息（推送者，可能与提交作者不同）
	PushUserID             *int      `json:"push_user_id,omitempty"`
	PushUserName           string    `json:"push_user_name,omitempty"`
	PushUserUsername       string    `json:"push_user_username,omitempty"`
	PushUserEmail          string    `json:"push_user_email,omitempty"`
	AddedFiles     []string          `json:"added_files"`
	ModifiedFiles  []string          `json:"modified_files"`
	RemovedFiles   []string          `json:"removed_files"`
	// FileStats 文件统计信息（可选，用于传递行数信息）
	FileStats map[string]*FileStat `json:"file_stats,omitempty"`
}

// FileStat 文件统计信息
type FileStat struct {
	AddedLines   int `json:"added_lines"`
	RemovedLines int `json:"removed_lines"`
}

