package commit

import (
	"fmt"
	"time"

	"gitlab-webhook-server/internal/model"
	"gitlab-webhook-server/internal/repository"
	"gitlab-webhook-server/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommitServiceV2 提交服务 V2（带数据库支持）
type CommitServiceV2 struct {
	logger   *zap.Logger
	repo     *repository.CommitRepository
	db       *gorm.DB
}

// NewCommitServiceV2 创建新的提交服务 V2
func NewCommitServiceV2(db *gorm.DB, logger *zap.Logger) *CommitServiceV2 {
	return &CommitServiceV2{
		logger: logger,
		repo:   repository.NewCommitRepository(db, logger),
		db:     db,
	}
}

// RecordCommit 记录代码提交（完整版本，包含行数和语言统计）
// commitRecord 可以包含 DiffStats 字段来传递行数信息
func (s *CommitServiceV2) RecordCommit(commitRecord *model.CommitRecord) error {
	// 检查是否已存在（使用 commit_id + project_id 唯一性）
	var existing model.Commit
	query := s.db.Where("commit_id = ?", commitRecord.CommitID)
	
	if commitRecord.ProjectID != nil {
		query = query.Where("project_id = ?", *commitRecord.ProjectID)
	} else {
		query = query.Where("project_id IS NULL")
	}
	
	err := query.First(&existing).Error
	if err == nil {
		s.logger.Info("提交记录已存在，跳过",
			zap.String("commit_id", commitRecord.CommitID),
			zap.Any("project_id", commitRecord.ProjectID),
		)
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询提交记录失败: %w", err)
	}

	// 解析时间戳
	timestamp, err := time.Parse(time.RFC3339, commitRecord.Timestamp)
	if err != nil {
		// 尝试其他格式
		timestamp, err = time.Parse("2006-01-02T15:04:05Z07:00", commitRecord.Timestamp)
		if err != nil {
			s.logger.Warn("解析时间戳失败，使用当前时间",
				zap.String("timestamp", commitRecord.Timestamp),
				zap.Error(err),
			)
			timestamp = time.Now()
		}
	}

	// 处理 authored_date 和 committed_date
	authoredDate := timestamp
	committedDate := timestamp
	if commitRecord.AuthoredDate != nil {
		authoredDate = *commitRecord.AuthoredDate
	}
	if commitRecord.CommittedDate != nil {
		committedDate = *commitRecord.CommittedDate
	}

	// 创建提交记录
	commit := &model.Commit{
		CommitID:              commitRecord.CommitID,
		ProjectID:             commitRecord.ProjectID,
		Message:               commitRecord.Message,
		Title:                 commitRecord.Title,
		Timestamp:             timestamp, // 保持向后兼容
		Author:                commitRecord.Author,
		AuthorEmail:           commitRecord.AuthorEmail,
		CommitterName:         commitRecord.CommitterName,
		CommitterEmail:        commitRecord.CommitterEmail,
		AuthoredDate:          &authoredDate,
		CommittedDate:         &committedDate,
		Branch:                commitRecord.Branch,
		RefProtected:          commitRecord.RefProtected,
		URL:                   commitRecord.URL,
		ProjectName:           commitRecord.ProjectName,
		ProjectPath:           commitRecord.ProjectPath,
		ProjectDescription:    commitRecord.ProjectDescription,
		ProjectWebURL:         commitRecord.ProjectWebURL,
		ProjectNamespace:      commitRecord.ProjectNamespace,
		ProjectVisibilityLevel: commitRecord.ProjectVisibilityLevel,
		ProjectDefaultBranch:  commitRecord.ProjectDefaultBranch,
		ProjectGitSSHURL:      commitRecord.ProjectGitSSHURL,
		ProjectGitHTTPURL:     commitRecord.ProjectGitHTTPURL,
		RepositoryName:        commitRecord.RepositoryName,
		RepositoryURL:         commitRecord.RepositoryURL,
		RepositoryDescription: commitRecord.RepositoryDescription,
		RepositoryHomepage:    commitRecord.RepositoryHomepage,
		RepositoryGitSSHURL:   commitRecord.RepositoryGitSSHURL,
		RepositoryGitHTTPURL:  commitRecord.RepositoryGitHTTPURL,
		RepositoryVisibilityLevel: commitRecord.RepositoryVisibilityLevel,
		BeforeSHA:             commitRecord.BeforeSHA,
		AfterSHA:              commitRecord.AfterSHA,
		CheckoutSHA:           commitRecord.CheckoutSHA,
		PushMessage:           commitRecord.PushMessage,
		TotalCommitsCount:     commitRecord.TotalCommitsCount,
		PushUserID:            commitRecord.PushUserID,
		PushUserName:          commitRecord.PushUserName,
		PushUserUsername:      commitRecord.PushUserUsername,
		PushUserEmail:         commitRecord.PushUserEmail,
		TotalAddedLines:       0,
		TotalRemovedLines:      0,
		TotalChangedFiles:      0,
	}

	// 处理文件变更
	var totalAdded, totalRemoved int
	languageStats := make(map[string]*LanguageFileStats)

	// 处理新增文件
	for _, filePath := range commitRecord.AddedFiles {
		addedLines, removedLines := s.getFileStats(commitRecord, filePath)
		file := s.createCommitFile(commit, filePath, "added", addedLines, removedLines)
		commit.Files = append(commit.Files, *file)
		totalAdded += file.AddedLines
		s.updateLanguageStats(languageStats, file.Language, file.AddedLines, 0, 1)
	}

	// 处理修改文件
	for _, filePath := range commitRecord.ModifiedFiles {
		addedLines, removedLines := s.getFileStats(commitRecord, filePath)
		file := s.createCommitFile(commit, filePath, "modified", addedLines, removedLines)
		commit.Files = append(commit.Files, *file)
		totalAdded += file.AddedLines
		totalRemoved += file.RemovedLines
		s.updateLanguageStats(languageStats, file.Language, file.AddedLines, file.RemovedLines, 1)
	}

	// 处理删除文件
	for _, filePath := range commitRecord.RemovedFiles {
		addedLines, removedLines := s.getFileStats(commitRecord, filePath)
		file := s.createCommitFile(commit, filePath, "removed", addedLines, removedLines)
		commit.Files = append(commit.Files, *file)
		totalRemoved += file.RemovedLines
		s.updateLanguageStats(languageStats, file.Language, 0, file.RemovedLines, 1)
	}

	// 更新总计
	commit.TotalAddedLines = totalAdded
	commit.TotalRemovedLines = totalRemoved
	commit.TotalChangedFiles = len(commit.Files)

	// 创建语言统计记录
	for lang, stats := range languageStats {
		commit.Languages = append(commit.Languages, model.CommitLanguage{
			Language:     lang,
			AddedLines:   stats.AddedLines,
			RemovedLines: stats.RemovedLines,
			FileCount:    stats.FileCount,
		})
	}

	// 使用事务保存
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(commit).Error; err != nil {
			return fmt.Errorf("保存提交记录失败: %w", err)
		}
		return nil
	}); err != nil {
		s.logger.Error("保存提交记录失败",
			zap.Error(err),
			zap.String("commit_id", commitRecord.CommitID),
		)
		return err
	}

	s.logger.Info("📝 提交记录已保存",
		zap.String("commit_id", commit.CommitID),
		zap.String("author", commit.Author),
		zap.String("project", commit.ProjectName),
		zap.Int("added_lines", commit.TotalAddedLines),
		zap.Int("removed_lines", commit.TotalRemovedLines),
		zap.Int("files", commit.TotalChangedFiles),
		zap.Int("languages", len(commit.Languages)),
	)

	return nil
}

// createCommitFile 创建文件变更记录
func (s *CommitServiceV2) createCommitFile(
	commit *model.Commit,
	filePath string,
	changeType string,
	addedLines, removedLines int,
) *model.CommitFile {
	language := utils.DetectLanguage(filePath)
	extension := utils.GetFileExtension(filePath)
	fileName := utils.GetFileName(filePath)

	return &model.CommitFile{
		FilePath:      filePath,
		FileName:      fileName,
		FileExtension: extension,
		ChangeType:    changeType,
		AddedLines:    addedLines,
		RemovedLines:  removedLines,
		Language:      language,
	}
}

// updateLanguageStats 更新语言统计
func (s *CommitServiceV2) updateLanguageStats(
	stats map[string]*LanguageFileStats,
	language string,
	added, removed int,
	fileCount int,
) {
	if stats[language] == nil {
		stats[language] = &LanguageFileStats{}
	}
	stats[language].AddedLines += added
	stats[language].RemovedLines += removed
	stats[language].FileCount += fileCount
}

// LanguageFileStats 语言文件统计
type LanguageFileStats struct {
	AddedLines   int
	RemovedLines int
	FileCount    int
}

// GetMemberCommits 获取成员的提交记录
func (s *CommitServiceV2) GetMemberCommits(
	authorEmail string,
	startDate, endDate *time.Time,
) ([]*model.Commit, error) {
	return s.repo.GetMemberCommits(authorEmail, startDate, endDate)
}

// GetMemberStats 获取成员统计信息
func (s *CommitServiceV2) GetMemberStats(
	authorEmail string,
	startDate, endDate *time.Time,
) (*repository.MemberStats, error) {
	return s.repo.GetMemberStats(authorEmail, startDate, endDate)
}

// GetLanguageStats 获取语言统计信息
func (s *CommitServiceV2) GetLanguageStats(
	authorEmail string,
	startDate, endDate *time.Time,
) ([]*repository.LanguageStats, error) {
	return s.repo.GetLanguageStats(authorEmail, startDate, endDate)
}

// getFileStats 获取文件统计信息
func (s *CommitServiceV2) getFileStats(commitRecord *model.CommitRecord, filePath string) (added, removed int) {
	if commitRecord.FileStats != nil {
		if stat, ok := commitRecord.FileStats[filePath]; ok {
			return stat.AddedLines, stat.RemovedLines
		}
	}
	return 0, 0
}

