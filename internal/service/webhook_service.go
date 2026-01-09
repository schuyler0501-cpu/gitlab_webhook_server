package service

import (
	"gitlab-webhook-server/internal/model"
	"gitlab-webhook-server/internal/queue"
	"gitlab-webhook-server/internal/service/commit"
	"gitlab-webhook-server/internal/webhook"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WebhookService Webhook 服务
type WebhookService struct {
	logger        *zap.Logger
	commitService *commit.CommitServiceV2
	db            *gorm.DB
	workerPool    *queue.WorkerPool
	webhookSecret string // Webhook 密钥（用于 token 验证）
}

// NewWebhookService 创建新的 Webhook 服务
func NewWebhookService(db *gorm.DB, workerPool *queue.WorkerPool, logger *zap.Logger) *WebhookService {
	return &WebhookService{
		logger:        logger,
		commitService: commit.NewCommitServiceV2(db, logger),
		db:            db,
		workerPool:    workerPool,
		webhookSecret: "", // 从配置中获取，需要在 handler 中设置
	}
}

// SetWebhookSecret 设置 webhook 密钥
func (s *WebhookService) SetWebhookSecret(secret string) {
	s.webhookSecret = secret
}

// GetWebhookSecret 获取 webhook 密钥
func (s *WebhookService) GetWebhookSecret() string {
	return s.webhookSecret
}

// ProcessWebhook 处理 webhook 事件
// platform: webhook 平台解析器
// eventType: 事件类型
// payload: webhook 负载数据
func (s *WebhookService) ProcessWebhook(platform webhook.Platform, eventType string, payload map[string]interface{}) error {
	s.logger.Info("收到 Webhook 事件",
		zap.String("platform", platform.GetPlatformName()),
		zap.String("event_type", eventType),
	)

	// 根据平台和事件类型处理
	switch eventType {
	case "Push Hook", "push": // GitLab/Gitee 使用 "Push Hook", GitHub 使用 "push"
		return s.handlePushEvent(platform, payload)
	case "Tag Push Hook", "tag_push": // GitLab/Gitee 使用 "Tag Push Hook"
		return s.handleTagPushEvent(platform, payload)
	default:
		s.logger.Info("未处理的事件类型",
			zap.String("platform", platform.GetPlatformName()),
			zap.String("event_type", eventType),
		)
		return nil
	}
}

// handlePushEvent 处理 Push 事件
func (s *WebhookService) handlePushEvent(platform webhook.Platform, payload map[string]interface{}) error {
	// 使用平台解析器解析提交记录
	commitRecords, err := platform.ParsePushEvent(payload)
	if err != nil {
		s.logger.Error("解析 Push 事件失败",
			zap.String("platform", platform.GetPlatformName()),
			zap.Error(err),
		)
		return err
	}

	if len(commitRecords) == 0 {
		s.logger.Info("Push 事件中没有提交记录",
			zap.String("platform", platform.GetPlatformName()),
		)
		return nil
	}

	// 异步处理提交记录
	if len(commitRecords) == 1 {
		// 单个提交，使用单任务
		task := queue.NewWebhookTask(commitRecords[0], s.commitService, s.logger)
		if err := s.workerPool.Submit(task); err != nil {
			s.logger.Error("提交任务失败", zap.Error(err))
			// 如果队列满，降级为同步处理
			if err := s.commitService.RecordCommit(commitRecords[0]); err != nil {
				s.logger.Error("记录提交失败", zap.Error(err))
			}
		}
	} else {
		// 批量提交，使用批量任务
		task := queue.NewBatchWebhookTask(commitRecords, s.commitService, s.db, s.logger)
		if err := s.workerPool.Submit(task); err != nil {
			s.logger.Error("提交批量任务失败", zap.Error(err))
			// 如果队列满，降级为同步处理
			for _, commitRecord := range commitRecords {
				if err := s.commitService.RecordCommit(commitRecord); err != nil {
					s.logger.Error("记录提交失败", zap.Error(err))
				}
			}
		}
	}

	return nil
}

// handleTagPushEvent 处理 Tag Push 事件
func (s *WebhookService) handleTagPushEvent(platform webhook.Platform, payload map[string]interface{}) error {
	s.logger.Info("处理 Tag Push 事件",
		zap.String("platform", platform.GetPlatformName()),
	)
	// TODO: 实现 Tag Push 事件处理逻辑
	return nil
}


