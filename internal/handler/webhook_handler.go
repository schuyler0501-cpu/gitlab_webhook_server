package handler

import (
	"net/http"
	"strings"

	"gitlab-webhook-server/internal/queue"
	"gitlab-webhook-server/internal/service"
	"gitlab-webhook-server/internal/webhook"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WebhookHandler Webhook 处理器
type WebhookHandler struct {
	logger         *zap.Logger
	webhookService *service.WebhookService
}

// NewWebhookHandler 创建新的 Webhook 处理器
func NewWebhookHandler(db *gorm.DB, workerPool *queue.WorkerPool, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		logger:         logger,
		webhookService: service.NewWebhookService(db, workerPool, logger),
	}
}

// HandleWebhook 处理 Webhook 请求（支持多平台：GitLab、Gitee、GitHub）
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// 收集所有请求头用于平台检测
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// 自动检测平台
	platform := webhook.DetectPlatform(headers)
	platformName := platform.GetPlatformName()
	eventType := platform.GetEventType(headers)

	// 获取 token（不同平台使用不同的 header）
	var token string
	switch platformName {
	case "gitlab":
		token = c.GetHeader("X-Gitlab-Token")
	case "gitee":
		token = c.GetHeader("X-Gitee-Token")
	case "github":
		// GitHub 使用 X-Hub-Signature-256 进行签名验证
		token = c.GetHeader("X-Hub-Signature-256")
	}

	// 解析请求体
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("解析 Webhook 请求失败",
			zap.String("platform", platformName),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 记录接收到的 webhook 信息（用于调试）
	h.logger.Debug("收到 Webhook 请求",
		zap.String("platform", platformName),
		zap.String("event_type", eventType),
		zap.Bool("has_token", token != ""),
	)

	// 异步处理 webhook（立即返回，不阻塞）
	go func() {
		if err := h.webhookService.ProcessWebhook(platform, eventType, payload); err != nil {
			h.logger.Error("处理 Webhook 失败",
				zap.String("platform", platformName),
				zap.Error(err),
			)
		}
	}()

	// 立即返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":  "Webhook received and queued for processing",
		"status":   "accepted",
		"platform": platformName,
		"event":    eventType,
	})
}

// Test Webhook 测试端点
func (h *WebhookHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook endpoint is ready",
	})
}

