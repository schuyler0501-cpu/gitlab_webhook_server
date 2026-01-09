package service

import (
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"
	gitlabClient "gitlab-webhook-server/internal/gitlab"
	"gitlab-webhook-server/internal/model"
	"gitlab-webhook-server/internal/service/commit"
	"gitlab-webhook-server/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ImportService 历史数据导入服务
type ImportService struct {
	logger        *zap.Logger
	gitlabClient  *gitlabClient.Client
	commitService *commit.CommitServiceV2
	db            *gorm.DB
}

// NewImportService 创建新的导入服务
func NewImportService(
	gitlabClient *gitlabClient.Client,
	commitService *commit.CommitServiceV2,
	db *gorm.DB,
	logger *zap.Logger,
) *ImportService {
	return &ImportService{
		logger:        logger,
		gitlabClient:  gitlabClient,
		commitService: commitService,
		db:            db,
	}
}

// ImportProjectCommits 导入项目的提交记录
func (s *ImportService) ImportProjectCommits(
	projectID string,
	since, until *time.Time,
	batchSize int,
) (*ImportResult, error) {
	result := &ImportResult{
		ProjectID: projectID,
		StartTime: time.Now(),
	}

	s.logger.Info("开始导入项目提交记录",
		zap.String("project_id", projectID),
		zap.Any("since", since),
		zap.Any("until", until),
	)

	// 获取项目信息
	project, _, err := s.gitlabClient.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("获取项目信息失败: %w", err)
	}

	projectName := project.Name
	projectPath := project.PathWithNamespace

	// 分页获取提交记录
	page := 1
	perPage := batchSize
	if perPage == 0 {
		perPage = 100 // 默认每页 100 条
	}

	for {
		commits, resp, err := s.gitlabClient.GetProjectCommits(projectID, since, until, page, perPage)
		if err != nil {
			return nil, fmt.Errorf("获取提交记录失败: %w", err)
		}

		if len(commits) == 0 {
			break
		}

		// 处理每批提交
		for _, commit := range commits {
			commitRecord, err := s.convertGitLabCommit(commit, projectName, projectPath, projectID)
			if err != nil {
				s.logger.Warn("转换提交记录失败",
					zap.String("commit_id", commit.ID),
					zap.Error(err),
				)
				result.Failed++
				continue
			}

			// 获取 diff 信息（包含行数统计）
			if diffs, err := s.gitlabClient.GetCommitDiff(projectID, commit.ID); err == nil {
				s.enrichCommitWithDiff(commitRecord, diffs)
			} else {
				s.logger.Debug("获取 diff 信息失败，将使用默认值",
					zap.String("commit_id", commit.ID),
					zap.Error(err),
				)
			}

			// 保存提交记录
			if err := s.commitService.RecordCommit(commitRecord); err != nil {
				s.logger.Warn("保存提交记录失败",
					zap.String("commit_id", commit.ID),
					zap.Error(err),
				)
				result.Failed++
				continue
			}

			result.Imported++
		}

		// 检查是否还有更多页
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage

		// 避免请求过快
		time.Sleep(time.Millisecond * 100)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	s.logger.Info("导入完成",
		zap.String("project_id", projectID),
		zap.Int("imported", result.Imported),
		zap.Int("failed", result.Failed),
		zap.Duration("duration", result.Duration),
	)

	return result, nil
}

// convertGitLabCommit 转换 GitLab Commit 为 CommitRecord
func (s *ImportService) convertGitLabCommit(
	commit *gitlab.Commit,
	projectName, projectPath, projectID string,
) (*model.CommitRecord, error) {
	// 解析时间
	timestamp := commit.CommittedDate.Format(time.RFC3339)

	// 获取作者信息
	authorName := "unknown"
	authorEmail := "unknown"
	if commit.AuthorName != "" {
		authorName = commit.AuthorName
	}
	if commit.AuthorEmail != "" {
		authorEmail = commit.AuthorEmail
	}

	// 解析文件变更（从 commit.Stats 或通过 diff 获取）
	addedFiles := make([]string, 0)
	modifiedFiles := make([]string, 0)
	removedFiles := make([]string, 0)

	return &model.CommitRecord{
		CommitID:      commit.ID,
		Message:       commit.Message,
		Timestamp:     timestamp,
		Author:        authorName,
		AuthorEmail:   authorEmail,
		URL:           commit.WebURL,
		ProjectName:   projectName,
		ProjectPath:   projectPath,
		AddedFiles:    addedFiles,
		ModifiedFiles: modifiedFiles,
		RemovedFiles:  removedFiles,
	}, nil
}

// enrichCommitWithDiff 使用 diff 信息丰富提交记录
func (s *ImportService) enrichCommitWithDiff(commitRecord *model.CommitRecord, diffs []*gitlab.Diff) {
	// 初始化 FileStats
	if commitRecord.FileStats == nil {
		commitRecord.FileStats = make(map[string]*model.FileStat)
	}

	for _, diff := range diffs {
		filePath := diff.NewPath
		if filePath == "" {
			filePath = diff.OldPath
		}

		// 从 diff 字符串中解析行数统计
		addedLines, removedLines := utils.ParseDiffStats(diff.Diff)

		// 记录文件统计信息
		commitRecord.FileStats[filePath] = &model.FileStat{
			AddedLines:   addedLines,
			RemovedLines: removedLines,
		}

		// 分类文件
		if diff.NewFile {
			commitRecord.AddedFiles = append(commitRecord.AddedFiles, filePath)
		} else if diff.DeletedFile {
			commitRecord.RemovedFiles = append(commitRecord.RemovedFiles, filePath)
		} else {
			commitRecord.ModifiedFiles = append(commitRecord.ModifiedFiles, filePath)
		}
	}
}

// ImportResult 导入结果
type ImportResult struct {
	ProjectID string
	Imported  int
	Failed    int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// GetImportStatus 获取导入状态
// 通过查询数据库中的提交记录来判断导入状态
func (s *ImportService) GetImportStatus(projectID string) (*ImportStatus, error) {
	status := &ImportStatus{
		ProjectID: projectID,
		Status:    "unknown",
	}

	// 查询数据库中该项目的提交记录数量
	var count int64
	query := s.db.Model(&model.Commit{}).Where("project_id = ?", projectID)
	if err := query.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询提交记录数量失败: %w", err)
	}

	status.TotalCommits = int(count)

	// 查询最近导入的记录时间
	var lastCommit model.Commit
	if err := query.Order("created_at DESC").First(&lastCommit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有记录，说明未导入
			status.Status = "not_started"
			status.Message = "该项目尚未导入任何提交记录"
			return status, nil
		}
		return nil, fmt.Errorf("查询最近提交记录失败: %w", err)
	}

	// 有记录，说明已导入
	status.Status = "completed"
	status.LastImportedAt = &lastCommit.CreatedAt
	status.Message = fmt.Sprintf("已导入 %d 条提交记录，最后导入时间: %s", count, lastCommit.CreatedAt.Format(time.RFC3339))

	return status, nil
}

// ImportStatus 导入状态
type ImportStatus struct {
	ProjectID      string     `json:"project_id"`
	Status         string     `json:"status"` // not_started, processing, completed, failed
	TotalCommits   int        `json:"total_commits"`
	LastImportedAt *time.Time `json:"last_imported_at,omitempty"`
	Message        string     `json:"message"`
}
