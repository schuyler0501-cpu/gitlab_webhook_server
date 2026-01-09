package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
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
	webhookSecret  string // Webhook 密钥（用于 token 验证）
}

// NewWebhookHandler 创建新的 Webhook 处理器
func NewWebhookHandler(db *gorm.DB, workerPool *queue.WorkerPool, webhookSecret string, logger *zap.Logger) *WebhookHandler {
	webhookService := service.NewWebhookService(db, workerPool, logger)
	webhookService.SetWebhookSecret(webhookSecret)
	return &WebhookHandler{
		logger:         logger,
		webhookService: webhookService,
		webhookSecret:  webhookSecret,
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

	// 读取请求体（GitHub 需要先读取以验证签名）
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("读取请求体失败",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// 恢复请求体供后续解析使用
	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

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

	// 验证 token（如果配置了 webhook secret）
	if h.webhookSecret != "" {
		if platformName == "github" {
			// GitHub 使用 HMAC SHA256 签名验证
			signature := c.GetHeader("X-Hub-Signature-256")
			if signature == "" {
				h.logger.Warn("GitHub webhook 签名缺失",
					zap.String("ip", c.ClientIP()),
				)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub webhook signature is required"})
				return
			}

			// 验证签名
			if !h.verifyGitHubSignature(bodyBytes, signature) {
				h.logger.Warn("GitHub webhook 签名验证失败",
					zap.String("ip", c.ClientIP()),
				)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid GitHub webhook signature"})
				return
			}
		} else {
			// GitLab/Gitee: 直接比较 token
			if token == "" {
				h.logger.Warn("Webhook token 缺失",
					zap.String("platform", platformName),
					zap.String("ip", c.ClientIP()),
				)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Webhook token is required"})
				return
			}
			if token != h.webhookSecret {
				h.logger.Warn("Webhook token 验证失败",
					zap.String("platform", platformName),
					zap.String("ip", c.ClientIP()),
				)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook token"})
				return
			}
		}
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

// verifyGitHubSignature 验证 GitHub webhook 签名
// GitHub 使用 HMAC SHA256 算法，签名格式为 "sha256=<hex_string>"
func (h *WebhookHandler) verifyGitHubSignature(payload []byte, signature string) bool {
	// 移除 "sha256=" 前缀
	signature = strings.TrimPrefix(signature, "sha256=")
	if signature == "" {
		return false
	}

	// 计算 HMAC SHA256
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 使用 hmac.Equal 进行常量时间比较，防止时序攻击
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Test Webhook 测试端点
func (h *WebhookHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook endpoint is ready",
	})
}

