package queue

import (
	"fmt"

	"gitlab-webhook-server/internal/model"
	"gitlab-webhook-server/internal/service/commit"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WebhookTask Webhook 处理任务
type WebhookTask struct {
	ID          string
	CommitRecord *model.CommitRecord
	commitService *commit.CommitServiceV2
	logger      *zap.Logger
}

// NewWebhookTask 创建新的 Webhook 任务
func NewWebhookTask(
	commitRecord *model.CommitRecord,
	commitService *commit.CommitServiceV2,
	logger *zap.Logger,
) *WebhookTask {
	return &WebhookTask{
		ID:           commitRecord.CommitID,
		CommitRecord: commitRecord,
		commitService: commitService,
		logger:       logger,
	}
}

// GetID 获取任务 ID
func (t *WebhookTask) GetID() string {
	return t.ID
}

// Execute 执行任务
func (t *WebhookTask) Execute() error {
	if err := t.commitService.RecordCommit(t.CommitRecord); err != nil {
		return fmt.Errorf("记录提交失败: %w", err)
	}
	return nil
}

// BatchWebhookTask 批量 Webhook 处理任务
type BatchWebhookTask struct {
	ID           string
	CommitRecords []*model.CommitRecord
	commitService *commit.CommitServiceV2
	db           *gorm.DB
	logger       *zap.Logger
}

// NewBatchWebhookTask 创建新的批量 Webhook 任务
func NewBatchWebhookTask(
	commitRecords []*model.CommitRecord,
	commitService *commit.CommitServiceV2,
	db *gorm.DB,
	logger *zap.Logger,
) *BatchWebhookTask {
	return &BatchWebhookTask{
		ID:            fmt.Sprintf("batch_%d", len(commitRecords)),
		CommitRecords: commitRecords,
		commitService: commitService,
		db:            db,
		logger:        logger,
	}
}

// GetID 获取任务 ID
func (t *BatchWebhookTask) GetID() string {
	return t.ID
}

// Execute 执行批量任务（使用事务）
func (t *BatchWebhookTask) Execute() error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		for _, commitRecord := range t.CommitRecords {
			if err := t.commitService.RecordCommit(commitRecord); err != nil {
				return fmt.Errorf("批量记录提交失败: %w", err)
			}
		}
		return nil
	})
}

