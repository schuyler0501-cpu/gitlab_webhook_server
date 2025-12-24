package handler

import (
	"net/http"

	"gitlab-webhook-server/internal/queue"
	"gitlab-webhook-server/internal/service"

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

// HandleWebhook 处理 GitLab Webhook 请求
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-Gitlab-Event")
	token := c.GetHeader("X-Gitlab-Token")

	// 验证 token（如果配置了）
	// 注意：这里需要从 config 中获取 secret，但为了保持 handler 的简洁性，
	// token 验证逻辑应该在 service 层或中间件中实现
	// 目前先记录 token 信息，后续可以在中间件中添加验证

	// 解析请求体
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("解析 Webhook 请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 记录接收到的 webhook 信息（用于调试）
	h.logger.Debug("收到 Webhook 请求",
		zap.String("event_type", eventType),
		zap.Bool("has_token", token != ""),
	)

	// 异步处理 webhook（立即返回，不阻塞）
	go func() {
		if err := h.webhookService.ProcessWebhook(eventType, payload); err != nil {
			h.logger.Error("处理 Webhook 失败", zap.Error(err))
		}
	}()

	// 立即返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook received and queued for processing",
		"status":  "accepted",
	})
}

// Test Webhook 测试端点
func (h *WebhookHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook endpoint is ready",
	})
}

