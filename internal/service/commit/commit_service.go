package commit

import (
	"fmt"

	"gitlab-webhook-server/internal/model"

	"go.uber.org/zap"
)

// CommitService 提交服务
type CommitService struct {
	logger *zap.Logger
}

// NewCommitService 创建新的提交服务
func NewCommitService(logger *zap.Logger) *CommitService {
	return &CommitService{
		logger: logger,
	}
}

// RecordCommit 记录代码提交
func (s *CommitService) RecordCommit(commit *model.CommitRecord) error {
	s.logger.Info("📝 记录代码提交",
		zap.String("commit_id", commit.CommitID),
		zap.String("author", commit.Author),
		zap.String("project", commit.ProjectName),
		zap.String("message", truncateString(commit.Message, 50)),
	)

	// TODO: 实现数据持久化逻辑
	// 这里可以：
	// 1. 保存到数据库（PostgreSQL, MySQL 等）
	// 2. 发送到消息队列（RabbitMQ, Kafka 等）
	// 3. 调用其他服务 API

	// 计算并记录统计信息
	stats := s.calculateCommitStats(commit)
	s.logger.Info("📊 提交统计",
		zap.Int("added_files", stats.AddedFiles),
		zap.Int("modified_files", stats.ModifiedFiles),
		zap.Int("removed_files", stats.RemovedFiles),
		zap.Int("total_changes", stats.TotalChanges),
	)

	return nil
}

// GetMemberCommits 获取成员的提交记录（用于统计）
func (s *CommitService) GetMemberCommits(
	authorEmail string,
	startDate, endDate *string,
) ([]*model.CommitRecord, error) {
	s.logger.Info("查询成员提交记录", zap.String("author_email", authorEmail))
	// TODO: 实现查询逻辑
	return []*model.CommitRecord{}, nil
}

// calculateCommitStats 计算提交统计信息
func (s *CommitService) calculateCommitStats(commit *model.CommitRecord) *CommitStats {
	return &CommitStats{
		AddedFiles:    len(commit.AddedFiles),
		ModifiedFiles: len(commit.ModifiedFiles),
		RemovedFiles:  len(commit.RemovedFiles),
		TotalChanges:  len(commit.AddedFiles) + len(commit.ModifiedFiles) + len(commit.RemovedFiles),
	}
}

// CommitStats 提交统计信息
type CommitStats struct {
	AddedFiles    int
	ModifiedFiles int
	RemovedFiles  int
	TotalChanges  int
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

