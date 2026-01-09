package commit

import (
	"gitlab-webhook-server/internal/model"

	"go.uber.org/zap"
)

// CommitService 提交服务（旧版本，已废弃）
// 注意：此服务已被 CommitServiceV2 替代，保留仅用于向后兼容
// 新代码请使用 CommitServiceV2，它包含完整的数据库持久化功能
type CommitService struct {
	logger *zap.Logger
}

// NewCommitService 创建新的提交服务（已废弃，请使用 NewCommitServiceV2）
func NewCommitService(logger *zap.Logger) *CommitService {
	return &CommitService{
		logger: logger,
	}
}

// RecordCommit 记录代码提交（已废弃，请使用 CommitServiceV2.RecordCommit）
// 此方法仅记录日志，不进行数据持久化
func (s *CommitService) RecordCommit(commit *model.CommitRecord) error {
	s.logger.Info("📝 记录代码提交",
		zap.String("commit_id", commit.CommitID),
		zap.String("author", commit.Author),
		zap.String("project", commit.ProjectName),
		zap.String("message", truncateString(commit.Message, 50)),
	)

	// 注意：此方法不进行数据持久化
	// 数据持久化功能已在 CommitServiceV2 中实现
	// 请使用 CommitServiceV2.RecordCommit 方法

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

// GetMemberCommits 获取成员的提交记录（已废弃，请使用 CommitServiceV2.GetMemberCommits）
// 此方法返回空列表，实际查询功能已在 CommitServiceV2 中实现
func (s *CommitService) GetMemberCommits(
	authorEmail string,
	startDate, endDate *string,
) ([]*model.CommitRecord, error) {
	s.logger.Info("查询成员提交记录", zap.String("author_email", authorEmail))
	// 注意：此方法不进行实际查询
	// 查询功能已在 CommitServiceV2 中实现
	// 请使用 CommitServiceV2.GetMemberCommits 方法
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

