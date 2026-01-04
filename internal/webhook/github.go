package webhook

import (
	"strings"
	"time"

	"gitlab-webhook-server/internal/model"
)

// GitHubPlatform GitHub 平台解析器
type GitHubPlatform struct{}

// NewGitHubPlatform 创建 GitHub 平台实例
func NewGitHubPlatform() *GitHubPlatform {
	return &GitHubPlatform{}
}

// GetPlatformName 获取平台名称
func (p *GitHubPlatform) GetPlatformName() string {
	return "github"
}

// Detect 检测是否为 GitHub webhook
func (p *GitHubPlatform) Detect(headers map[string]string) bool {
	eventHeader := headers["X-GitHub-Event"]
	return eventHeader != ""
}

// GetEventType 获取事件类型
func (p *GitHubPlatform) GetEventType(headers map[string]string) string {
	return headers["X-GitHub-Event"]
}

// ParsePushEvent 解析 GitHub Push 事件
func (p *GitHubPlatform) ParsePushEvent(payload map[string]interface{}) ([]*model.CommitRecord, error) {
	commits, ok := payload["commits"].([]interface{})
	if !ok || len(commits) == 0 {
		return []*model.CommitRecord{}, nil
	}

	// 解析推送级别信息（所有提交共享）
	pushInfo := p.parsePushInfo(payload)

	// 收集所有提交记录
	var commitRecords []*model.CommitRecord
	for _, commitData := range commits {
		commitMap, ok := commitData.(map[string]interface{})
		if !ok {
			continue
		}

		commitRecord := p.parseCommit(commitMap, pushInfo)
		if commitRecord != nil {
			commitRecords = append(commitRecords, commitRecord)
		}
	}

	return commitRecords, nil
}

// PushInfo 推送级别信息（所有提交共享）
type GitHubPushInfo struct {
	ProjectName              string
	ProjectPath              string
	ProjectID                *int
	ProjectDescription       string
	ProjectWebURL            string
	ProjectNamespace         string
	ProjectVisibilityLevel   *int
	ProjectDefaultBranch     string
	ProjectGitSSHURL         string
	ProjectGitHTTPURL        string
	RepositoryName           string
	RepositoryURL            string
	RepositoryDescription    string
	RepositoryHomepage       string
	RepositoryGitSSHURL      string
	RepositoryGitHTTPURL     string
	RepositoryVisibilityLevel *int
	Branch                   string
	RefProtected             *bool
	BeforeSHA                string
	AfterSHA                 string
	CheckoutSHA              string
	PushMessage              string
	TotalCommitsCount        int
	PushUserID               *int
	PushUserName              string
	PushUserUsername          string
	PushUserEmail             string
}

// parsePushInfo 解析推送级别信息
func (p *GitHubPlatform) parsePushInfo(payload map[string]interface{}) *GitHubPushInfo {
	info := &GitHubPushInfo{}

	// 解析仓库信息（GitHub 使用 repository 字段）
	repository, _ := payload["repository"].(map[string]interface{})
	if repository != nil {
		if name, ok := repository["name"].(string); ok {
			info.ProjectName = name
			info.RepositoryName = name
		}
		// GitHub 使用 full_name
		if fullName, ok := repository["full_name"].(string); ok {
			info.ProjectPath = fullName
		}
		if repoID, ok := repository["id"].(float64); ok {
			id := int(repoID)
			info.ProjectID = &id
		}
		if desc, ok := repository["description"].(string); ok {
			info.ProjectDescription = desc
			info.RepositoryDescription = desc
		}
		if htmlURL, ok := repository["html_url"].(string); ok {
			info.ProjectWebURL = htmlURL
			info.RepositoryURL = htmlURL
			info.RepositoryHomepage = htmlURL
		}
		// GitHub 使用 owner
		if owner, ok := repository["owner"].(map[string]interface{}); ok {
			if login, ok := owner["login"].(string); ok {
				info.ProjectNamespace = login
			}
		}
		// GitHub 使用 private 字段
		if private, ok := repository["private"].(bool); ok {
			level := 0
			if !private {
				level = 20 // 公开
			} else {
				level = 0 // 私有
			}
			info.ProjectVisibilityLevel = &level
			info.RepositoryVisibilityLevel = &level
		}
		if defaultBranch, ok := repository["default_branch"].(string); ok {
			info.ProjectDefaultBranch = defaultBranch
		}
		if gitSSH, ok := repository["ssh_url"].(string); ok {
			info.ProjectGitSSHURL = gitSSH
			info.RepositoryGitSSHURL = gitSSH
		}
		if gitHTTP, ok := repository["clone_url"].(string); ok {
			info.ProjectGitHTTPURL = gitHTTP
			info.RepositoryGitHTTPURL = gitHTTP
		}
	}

	// 解析分支信息
	if ref, ok := payload["ref"].(string); ok {
		if strings.HasPrefix(ref, "refs/heads/") {
			info.Branch = strings.TrimPrefix(ref, "refs/heads/")
		} else {
			info.Branch = ref
		}
	}

	// 解析分支保护状态（GitHub 在 repository 中）
	if repository != nil {
		if protected, ok := repository["protected"].(bool); ok {
			info.RefProtected = &protected
		}
	}

	// 解析推送 SHA
	if before, ok := payload["before"].(string); ok {
		info.BeforeSHA = before
	}
	if after, ok := payload["after"].(string); ok {
		info.AfterSHA = after
		info.CheckoutSHA = after
	}

	// 解析总提交数
	if commits, ok := payload["commits"].([]interface{}); ok {
		info.TotalCommitsCount = len(commits)
	}

	// 解析推送用户信息（GitHub 使用 pusher）
	if pusher, ok := payload["pusher"].(map[string]interface{}); ok {
		if name, ok := pusher["name"].(string); ok {
			info.PushUserName = name
			info.PushUserUsername = name
		}
		if email, ok := pusher["email"].(string); ok {
			info.PushUserEmail = email
		}
	} else if sender, ok := payload["sender"].(map[string]interface{}); ok {
		if login, ok := sender["login"].(string); ok {
			info.PushUserUsername = login
		}
		if name, ok := sender["name"].(string); ok {
			info.PushUserName = name
		}
		if email, ok := sender["email"].(string); ok {
			info.PushUserEmail = email
		}
	}

	return info
}

// parseCommit 解析提交数据
func (p *GitHubPlatform) parseCommit(
	commitMap map[string]interface{},
	pushInfo *GitHubPushInfo,
) *model.CommitRecord {
	// GitHub 使用 id 或 sha
	commitID, _ := commitMap["id"].(string)
	if commitID == "" {
		commitID, _ = commitMap["sha"].(string)
	}
	message, _ := commitMap["message"].(string)
	timestamp, _ := commitMap["timestamp"].(string)
	url, _ := commitMap["url"].(string)

	// GitHub 的 author 结构
	author, _ := commitMap["author"].(map[string]interface{})
	authorName := "unknown"
	authorEmail := "unknown"
	if author != nil {
		authorName, _ = author["name"].(string)
		authorEmail, _ = author["email"].(string)
	}

	// GitHub 的文件变更信息在 modified, added, removed 字段
	addedFiles := p.parseStringSlice(commitMap["added"])
	modifiedFiles := p.parseStringSlice(commitMap["modified"])
	removedFiles := p.parseStringSlice(commitMap["removed"])

	if commitID == "" {
		return nil
	}

	// 提取提交标题（message 第一行）
	title := message
	if newlineIdx := strings.Index(message, "\n"); newlineIdx > 0 {
		title = message[:newlineIdx]
	}
	if len(title) > 255 {
		title = title[:255]
	}

	// 解析时间戳（GitHub 使用 ISO 8601 格式）
	var authoredDate, committedDate *time.Time
	if timestamp != "" {
		// 尝试 RFC3339 格式
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
			committedDate = &t
			authoredDate = &t
		} else {
			// 尝试其他常见格式
			formats := []string{
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02T15:04:05Z",
			}
			for _, format := range formats {
				if t, err := time.Parse(format, timestamp); err == nil {
					committedDate = &t
					authoredDate = &t
					break
				}
			}
		}
	}

	// 获取提交者信息（可能与作者不同）
	committerName := authorName
	committerEmail := authorEmail
	if committer, ok := commitMap["committer"].(map[string]interface{}); ok {
		if name, ok := committer["name"].(string); ok {
			committerName = name
		}
		if email, ok := committer["email"].(string); ok {
			committerEmail = email
		}
	}

	return &model.CommitRecord{
		CommitID:                 commitID,
		Message:                  message,
		Title:                    title,
		Timestamp:                timestamp,
		Author:                   authorName,
		AuthorEmail:              authorEmail,
		CommitterName:            committerName,
		CommitterEmail:           committerEmail,
		AuthoredDate:             authoredDate,
		CommittedDate:            committedDate,
		Branch:                   pushInfo.Branch,
		RefProtected:             pushInfo.RefProtected,
		ProjectID:                pushInfo.ProjectID,
		URL:                      url,
		ProjectName:              pushInfo.ProjectName,
		ProjectPath:              pushInfo.ProjectPath,
		ProjectDescription:       pushInfo.ProjectDescription,
		ProjectWebURL:            pushInfo.ProjectWebURL,
		ProjectNamespace:         pushInfo.ProjectNamespace,
		ProjectVisibilityLevel:  pushInfo.ProjectVisibilityLevel,
		ProjectDefaultBranch:     pushInfo.ProjectDefaultBranch,
		ProjectGitSSHURL:         pushInfo.ProjectGitSSHURL,
		ProjectGitHTTPURL:        pushInfo.ProjectGitHTTPURL,
		RepositoryName:           pushInfo.RepositoryName,
		RepositoryURL:            pushInfo.RepositoryURL,
		RepositoryDescription:    pushInfo.RepositoryDescription,
		RepositoryHomepage:      pushInfo.RepositoryHomepage,
		RepositoryGitSSHURL:      pushInfo.RepositoryGitSSHURL,
		RepositoryGitHTTPURL:     pushInfo.RepositoryGitHTTPURL,
		RepositoryVisibilityLevel: pushInfo.RepositoryVisibilityLevel,
		BeforeSHA:                pushInfo.BeforeSHA,
		AfterSHA:                 pushInfo.AfterSHA,
		CheckoutSHA:              pushInfo.CheckoutSHA,
		PushMessage:              pushInfo.PushMessage,
		TotalCommitsCount:        pushInfo.TotalCommitsCount,
		PushUserID:               pushInfo.PushUserID,
		PushUserName:             pushInfo.PushUserName,
		PushUserUsername:         pushInfo.PushUserUsername,
		PushUserEmail:            pushInfo.PushUserEmail,
		AddedFiles:               addedFiles,
		ModifiedFiles:            modifiedFiles,
		RemovedFiles:             removedFiles,
	}
}

// parseStringSlice 解析字符串切片
func (p *GitHubPlatform) parseStringSlice(data interface{}) []string {
	if data == nil {
		return []string{}
	}

	slice, ok := data.([]interface{})
	if !ok {
		return []string{}
	}

	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

