package repository

import (
	"fmt"
	"time"

	"gitlab-webhook-server/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommitRepository 提交记录仓库
type CommitRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewCommitRepository 创建新的提交记录仓库
func NewCommitRepository(db *gorm.DB, logger *zap.Logger) *CommitRepository {
	return &CommitRepository{
		db:     db,
		logger: logger,
	}
}

// CreateCommit 创建提交记录
func (r *CommitRepository) CreateCommit(commit *model.Commit) error {
	if err := r.db.Create(commit).Error; err != nil {
		r.logger.Error("创建提交记录失败",
			zap.Error(err),
			zap.String("commit_id", commit.CommitID),
		)
		return fmt.Errorf("创建提交记录失败: %w", err)
	}

	r.logger.Info("提交记录创建成功",
		zap.Uint64("id", commit.ID),
		zap.String("commit_id", commit.CommitID),
	)

	return nil
}

// GetCommitByCommitID 根据提交 ID 获取提交记录
func (r *CommitRepository) GetCommitByCommitID(commitID string) (*model.Commit, error) {
	var commit model.Commit
	if err := r.db.Where("commit_id = ?", commitID).First(&commit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询提交记录失败: %w", err)
	}
	return &commit, nil
}

// GetMemberCommits 获取成员的提交记录
func (r *CommitRepository) GetMemberCommits(
	authorEmail string,
	startDate, endDate *time.Time,
) ([]*model.Commit, error) {
	var commits []*model.Commit
	query := r.db.Where("author_email = ?", authorEmail)

	if startDate != nil {
		query = query.Where("timestamp >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("timestamp < ?", *endDate)
	}

	if err := query.Order("timestamp DESC").Find(&commits).Error; err != nil {
		return nil, fmt.Errorf("查询成员提交记录失败: %w", err)
	}

	return commits, nil
}

// GetMemberStats 获取成员统计信息
func (r *CommitRepository) GetMemberStats(
	authorEmail string,
	startDate, endDate *time.Time,
) (*MemberStats, error) {
	var stats MemberStats
	query := r.db.Model(&model.Commit{}).Where("author_email = ?", authorEmail)

	if startDate != nil {
		query = query.Where("timestamp >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("timestamp < ?", *endDate)
	}

	if err := query.Select(
		"COUNT(*) as commit_count",
		"COALESCE(SUM(total_added_lines), 0) as total_added",
		"COALESCE(SUM(total_removed_lines), 0) as total_removed",
		"COALESCE(SUM(total_changed_files), 0) as total_files",
	).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("查询成员统计失败: %w", err)
	}

	return &stats, nil
}

// GetLanguageStats 获取语言统计信息
func (r *CommitRepository) GetLanguageStats(
	authorEmail string,
	startDate, endDate *time.Time,
) ([]*LanguageStats, error) {
	var stats []*LanguageStats

	query := r.db.Table("commit_languages").
		Select(
			"commit_languages.language",
			"COALESCE(SUM(commit_languages.added_lines), 0) as total_added",
			"COALESCE(SUM(commit_languages.removed_lines), 0) as total_removed",
			"COALESCE(SUM(commit_languages.file_count), 0) as total_files",
		).
		Joins("JOIN commits ON commit_languages.commit_id = commits.id").
		Where("commits.author_email = ?", authorEmail)

	if startDate != nil {
		query = query.Where("commits.timestamp >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("commits.timestamp < ?", *endDate)
	}

	if err := query.Group("commit_languages.language").
		Order("total_added DESC").
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("查询语言统计失败: %w", err)
	}

	return stats, nil
}

// MemberStats 成员统计信息
type MemberStats struct {
	CommitCount int `json:"commit_count"`
	TotalAdded  int `json:"total_added"`
	TotalRemoved int `json:"total_removed"`
	TotalFiles  int `json:"total_files"`
}

// LanguageStats 语言统计信息
type LanguageStats struct {
	Language    string `json:"language"`
	TotalAdded  int    `json:"total_added"`
	TotalRemoved int   `json:"total_removed"`
	TotalFiles  int    `json:"total_files"`
}

