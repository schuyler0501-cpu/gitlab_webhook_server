package webhook

import (
	"strings"
	"time"

	"gitlab-webhook-server/internal/model"
)

// GitLabPlatform GitLab 平台解析器
type GitLabPlatform struct{}

// NewGitLabPlatform 创建 GitLab 平台实例
func NewGitLabPlatform() *GitLabPlatform {
	return &GitLabPlatform{}
}

// GetPlatformName 获取平台名称
func (p *GitLabPlatform) GetPlatformName() string {
	return "gitlab"
}

// Detect 检测是否为 GitLab webhook
func (p *GitLabPlatform) Detect(headers map[string]string) bool {
	eventHeader := headers["X-Gitlab-Event"]
	return eventHeader != ""
}

// GetEventType 获取事件类型
func (p *GitLabPlatform) GetEventType(headers map[string]string) string {
	return headers["X-Gitlab-Event"]
}

// ParsePushEvent 解析 GitLab Push 事件
func (p *GitLabPlatform) ParsePushEvent(payload map[string]interface{}) ([]*model.CommitRecord, error) {
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
type PushInfo struct {
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
func (p *GitLabPlatform) parsePushInfo(payload map[string]interface{}) *PushInfo {
	info := &PushInfo{}

	// 解析项目信息
	project, _ := payload["project"].(map[string]interface{})
	if project != nil {
		if name, ok := project["name"].(string); ok {
			info.ProjectName = name
		}
		if path, ok := project["path_with_namespace"].(string); ok {
			info.ProjectPath = path
		}
		if pid, ok := project["id"].(float64); ok {
			id := int(pid)
			info.ProjectID = &id
		}
		if desc, ok := project["description"].(string); ok {
			info.ProjectDescription = desc
		}
		if webURL, ok := project["web_url"].(string); ok {
			info.ProjectWebURL = webURL
		}
		if namespace, ok := project["namespace"].(string); ok {
			info.ProjectNamespace = namespace
		}
		if visibility, ok := project["visibility_level"].(float64); ok {
			level := int(visibility)
			info.ProjectVisibilityLevel = &level
		}
		if defaultBranch, ok := project["default_branch"].(string); ok {
			info.ProjectDefaultBranch = defaultBranch
		}
		if gitSSH, ok := project["git_ssh_url"].(string); ok {
			info.ProjectGitSSHURL = gitSSH
		}
		if gitHTTP, ok := project["git_http_url"].(string); ok {
			info.ProjectGitHTTPURL = gitHTTP
		}
	}

	// 解析仓库信息
	repository, _ := payload["repository"].(map[string]interface{})
	if repository != nil {
		if name, ok := repository["name"].(string); ok {
			info.RepositoryName = name
		}
		if url, ok := repository["url"].(string); ok {
			info.RepositoryURL = url
		}
		if desc, ok := repository["description"].(string); ok {
			info.RepositoryDescription = desc
		}
		if homepage, ok := repository["homepage"].(string); ok {
			info.RepositoryHomepage = homepage
		}
		if gitSSH, ok := repository["git_ssh_url"].(string); ok {
			info.RepositoryGitSSHURL = gitSSH
		}
		if gitHTTP, ok := repository["git_http_url"].(string); ok {
			info.RepositoryGitHTTPURL = gitHTTP
		}
		if visibility, ok := repository["visibility_level"].(float64); ok {
			level := int(visibility)
			info.RepositoryVisibilityLevel = &level
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

	// 解析分支保护状态
	if protected, ok := payload["ref_protected"].(bool); ok {
		info.RefProtected = &protected
	}

	// 解析推送 SHA
	if before, ok := payload["before"].(string); ok {
		info.BeforeSHA = before
	}
	if after, ok := payload["after"].(string); ok {
		info.AfterSHA = after
	}
	if checkout, ok := payload["checkout_sha"].(string); ok {
		info.CheckoutSHA = checkout
	}

	// 解析推送消息
	if msg, ok := payload["message"].(string); ok {
		info.PushMessage = msg
	}

	// 解析总提交数
	if count, ok := payload["total_commits_count"].(float64); ok {
		info.TotalCommitsCount = int(count)
	}

	// 解析推送用户信息
	if userID, ok := payload["user_id"].(float64); ok {
		id := int(userID)
		info.PushUserID = &id
	}
	if userName, ok := payload["user_name"].(string); ok {
		info.PushUserName = userName
	}
	if userUsername, ok := payload["user_username"].(string); ok {
		info.PushUserUsername = userUsername
	}
	if userEmail, ok := payload["user_email"].(string); ok {
		info.PushUserEmail = userEmail
	}

	return info
}

// parseCommit 解析提交数据
func (p *GitLabPlatform) parseCommit(
	commitMap map[string]interface{},
	pushInfo *PushInfo,
) *model.CommitRecord {
	commitID, _ := commitMap["id"].(string)
	message, _ := commitMap["message"].(string)
	timestamp, _ := commitMap["timestamp"].(string)
	url, _ := commitMap["url"].(string)

	author, _ := commitMap["author"].(map[string]interface{})
	authorName := "unknown"
	authorEmail := "unknown"
	if author != nil {
		authorName, _ = author["name"].(string)
		authorEmail, _ = author["email"].(string)
	}

	// 解析文件变更
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

	// 解析时间戳
	var authoredDate, committedDate *time.Time
	if timestamp != "" {
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
			committedDate = &t
			authoredDate = &t // 默认相同，如果有区分可以单独解析
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
func (p *GitLabPlatform) parseStringSlice(data interface{}) []string {
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

