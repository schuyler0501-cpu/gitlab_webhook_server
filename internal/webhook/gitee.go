package webhook

import (
	"strings"
	"time"

	"gitlab-webhook-server/internal/model"
)

// GiteePlatform Gitee 平台解析器
type GiteePlatform struct{}

// NewGiteePlatform 创建 Gitee 平台实例
func NewGiteePlatform() *GiteePlatform {
	return &GiteePlatform{}
}

// GetPlatformName 获取平台名称
func (p *GiteePlatform) GetPlatformName() string {
	return "gitee"
}

// Detect 检测是否为 Gitee webhook
func (p *GiteePlatform) Detect(headers map[string]string) bool {
	eventHeader := headers["X-Gitee-Event"]
	return eventHeader != ""
}

// GetEventType 获取事件类型
func (p *GiteePlatform) GetEventType(headers map[string]string) string {
	return headers["X-Gitee-Event"]
}

// ParsePushEvent 解析 Gitee Push 事件
func (p *GiteePlatform) ParsePushEvent(payload map[string]interface{}) ([]*model.CommitRecord, error) {
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
type GiteePushInfo struct {
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
func (p *GiteePlatform) parsePushInfo(payload map[string]interface{}) *GiteePushInfo {
	info := &GiteePushInfo{}

	// 解析项目信息（Gitee 使用 project 字段）
	project, _ := payload["project"].(map[string]interface{})
	if project != nil {
		if name, ok := project["name"].(string); ok {
			info.ProjectName = name
		}
		// Gitee 使用 path_with_namespace 或 full_name
		if path, ok := project["path_with_namespace"].(string); ok {
			info.ProjectPath = path
		} else if fullName, ok := project["full_name"].(string); ok {
			info.ProjectPath = fullName
		}
		if pid, ok := project["id"].(float64); ok {
			id := int(pid)
			info.ProjectID = &id
		}
		if desc, ok := project["description"].(string); ok {
			info.ProjectDescription = desc
		}
		if webURL, ok := project["html_url"].(string); ok {
			info.ProjectWebURL = webURL
		} else if webURL, ok := project["url"].(string); ok {
			info.ProjectWebURL = webURL
		}
		// Gitee 使用 namespace 或 owner
		if namespace, ok := project["namespace"].(string); ok {
			info.ProjectNamespace = namespace
		} else if owner, ok := project["owner"].(map[string]interface{}); ok {
			if login, ok := owner["login"].(string); ok {
				info.ProjectNamespace = login
			}
		}
		// Gitee 使用 public/private 或 visibility
		if public, ok := project["public"].(bool); ok {
			level := 0
			if public {
				level = 20 // 公开
			} else {
				level = 0 // 私有
			}
			info.ProjectVisibilityLevel = &level
		}
		if defaultBranch, ok := project["default_branch"].(string); ok {
			info.ProjectDefaultBranch = defaultBranch
		}
		if gitSSH, ok := project["ssh_url"].(string); ok {
			info.ProjectGitSSHURL = gitSSH
		}
		if gitHTTP, ok := project["clone_url"].(string); ok {
			info.ProjectGitHTTPURL = gitHTTP
		} else if gitHTTP, ok := project["git_http_url"].(string); ok {
			info.ProjectGitHTTPURL = gitHTTP
		}
	}

	// Gitee 的 repository 信息通常在 project 中，单独处理
	info.RepositoryName = info.ProjectName
	info.RepositoryURL = info.ProjectWebURL
	info.RepositoryDescription = info.ProjectDescription
	info.RepositoryHomepage = info.ProjectWebURL
	info.RepositoryGitSSHURL = info.ProjectGitSSHURL
	info.RepositoryGitHTTPURL = info.ProjectGitHTTPURL
	info.RepositoryVisibilityLevel = info.ProjectVisibilityLevel

	// 解析分支信息
	if ref, ok := payload["ref"].(string); ok {
		if strings.HasPrefix(ref, "refs/heads/") {
			info.Branch = strings.TrimPrefix(ref, "refs/heads/")
		} else {
			info.Branch = ref
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
	if count, ok := payload["total_commits_count"].(float64); ok {
		info.TotalCommitsCount = int(count)
	} else if commits, ok := payload["commits"].([]interface{}); ok {
		info.TotalCommitsCount = len(commits)
	}

	// 解析推送用户信息（Gitee 使用 pusher 或 user）
	if pusher, ok := payload["pusher"].(map[string]interface{}); ok {
		if name, ok := pusher["name"].(string); ok {
			info.PushUserName = name
			info.PushUserUsername = name
		}
		if email, ok := pusher["email"].(string); ok {
			info.PushUserEmail = email
		}
	} else if user, ok := payload["user"].(map[string]interface{}); ok {
		if name, ok := user["name"].(string); ok {
			info.PushUserName = name
		}
		if login, ok := user["login"].(string); ok {
			info.PushUserUsername = login
		}
		if email, ok := user["email"].(string); ok {
			info.PushUserEmail = email
		}
	}

	return info
}

// parseCommit 解析提交数据
func (p *GiteePlatform) parseCommit(
	commitMap map[string]interface{},
	pushInfo *GiteePushInfo,
) *model.CommitRecord {
	// Gitee 使用 id 或 sha
	commitID, _ := commitMap["id"].(string)
	if commitID == "" {
		commitID, _ = commitMap["sha"].(string)
	}
	message, _ := commitMap["message"].(string)
	timestamp, _ := commitMap["timestamp"].(string)
	url, _ := commitMap["url"].(string)

	// Gitee 的 author 结构
	author, _ := commitMap["author"].(map[string]interface{})
	authorName := "unknown"
	authorEmail := "unknown"
	if author != nil {
		authorName, _ = author["name"].(string)
		authorEmail, _ = author["email"].(string)
	}

	// Gitee 的文件变更信息（可能在不同字段）
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

	// 解析时间戳（Gitee 可能使用不同格式）
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
				"2006-01-02 15:04:05",
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
func (p *GiteePlatform) parseStringSlice(data interface{}) []string {
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

// ParseTagPushEvent 解析 Gitee Tag Push 事件
// Tag Push 事件与 Push 事件结构相同，只是 ref 是 "refs/tags/" 开头
func (p *GiteePlatform) ParseTagPushEvent(payload map[string]interface{}) ([]*model.CommitRecord, error) {
	// Tag Push 事件结构与 Push 事件相同，可以直接复用 ParsePushEvent
	return p.ParsePushEvent(payload)
}

// VerifySecret 验证 Gitee webhook 密钥
func (p *GiteePlatform) VerifySecret(headers map[string]string, payload []byte, secret string) error {
	// Gitee 使用简单的 token 比较，在 handler 中已处理
	return nil
}
