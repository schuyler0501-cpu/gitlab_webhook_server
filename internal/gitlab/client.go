package gitlab

import (
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

// Client GitLab API 客户端
type Client struct {
	client *gitlab.Client
	logger *zap.Logger
}

// NewClient 创建新的 GitLab 客户端
func NewClient(baseURL, token string, logger *zap.Logger) (*Client, error) {
	var client *gitlab.Client
	var err error

	if baseURL != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	} else {
		client, err = gitlab.NewClient(token)
	}

	if err != nil {
		return nil, fmt.Errorf("创建 GitLab 客户端失败: %w", err)
	}

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

// GetProjectCommits 获取项目的提交记录
func (c *Client) GetProjectCommits(
	projectID string,
	since, until *time.Time,
	page, perPage int,
) ([]*gitlab.Commit, *gitlab.Response, error) {
	opts := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	if since != nil {
		opts.Since = since
	}
	if until != nil {
		opts.Until = until
	}

	commits, resp, err := c.client.Commits.ListCommits(projectID, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("获取提交记录失败: %w", err)
	}

	return commits, resp, nil
}

// GetCommitDiff 获取提交的 diff 信息（包含行数统计）
func (c *Client) GetCommitDiff(projectID, sha string) ([]*gitlab.Diff, error) {
	diffs, _, err := c.client.Commits.GetCommitDiff(projectID, sha, nil)
	if err != nil {
		return nil, fmt.Errorf("获取提交 diff 失败: %w", err)
	}

	return diffs, nil
}

// GetProject 获取项目信息
func (c *Client) GetProject(projectID string) (*gitlab.Project, *gitlab.Response, error) {
	project, resp, err := c.client.Projects.GetProject(projectID, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("获取项目信息失败: %w", err)
	}

	return project, resp, nil
}

// CalculateDiffStats 计算 diff 统计信息
func CalculateDiffStats(diffs []*gitlab.Diff) (added, removed int) {
	for _, diff := range diffs {
		added += diff.Additions
		removed += diff.Deletions
	}
	return added, removed
}

