package webhook

import (
	"gitlab-webhook-server/internal/model"
)

// Platform Webhook 平台接口
// 定义统一的 webhook 解析接口，支持不同 Git 平台
type Platform interface {
	// Detect 检测是否为该平台的 webhook
	// 通过请求头信息判断
	Detect(headers map[string]string) bool

	// ParsePushEvent 解析 Push 事件
	// 返回提交记录列表和错误
	ParsePushEvent(payload map[string]interface{}) ([]*model.CommitRecord, error)

	// GetEventType 获取事件类型
	// 从请求头中提取事件类型
	GetEventType(headers map[string]string) string

	// GetPlatformName 获取平台名称
	GetPlatformName() string
}

// PlatformType 平台类型
type PlatformType string

const (
	PlatformGitLab PlatformType = "gitlab"
	PlatformGitee  PlatformType = "gitee"
	PlatformGitHub PlatformType = "github"
)

// GetPlatform 根据平台类型获取平台实例
func GetPlatform(platformType PlatformType) Platform {
	switch platformType {
	case PlatformGitLab:
		return NewGitLabPlatform()
	case PlatformGitee:
		return NewGiteePlatform()
	case PlatformGitHub:
		return NewGitHubPlatform()
	default:
		return NewGitLabPlatform() // 默认使用 GitLab
	}
}

// DetectPlatform 自动检测平台类型
// 根据请求头信息自动识别平台
func DetectPlatform(headers map[string]string) Platform {
	// 按优先级检测：GitLab -> Gitee -> GitHub
	platforms := []Platform{
		NewGitLabPlatform(),
		NewGiteePlatform(),
		NewGitHubPlatform(),
	}

	for _, platform := range platforms {
		if platform.Detect(headers) {
			return platform
		}
	}

	// 默认返回 GitLab（向后兼容）
	return NewGitLabPlatform()
}

