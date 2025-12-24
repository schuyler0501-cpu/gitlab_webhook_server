package router

import (
	"time"

	"gitlab-webhook-server/internal/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// New 创建新的路由实例
func New(logger *zap.Logger) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(ginLogger(logger))
	r.Use(gin.Recovery())

	return r
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(
	r *gin.Engine,
	webhookHandler *handler.WebhookHandler,
	statsHandler *handler.StatsHandler,
	importHandler *handler.ImportHandler,
) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Webhook 路由组
	webhook := r.Group("/webhook")
	{
		webhook.POST("", webhookHandler.HandleWebhook)
		webhook.GET("/test", webhookHandler.Test)
	}

	// 统计 API 路由组
	api := r.Group("/api/stats")
	{
		api.GET("/member", statsHandler.GetMemberStats)
		api.GET("/languages", statsHandler.GetLanguageStats)
		api.GET("/commits", statsHandler.GetMemberCommits)
	}

	// 导入 API 路由组（仅在 importHandler 不为 nil 时注册）
	if importHandler != nil {
		importAPI := r.Group("/api/import")
		{
			importAPI.POST("/project", importHandler.ImportProject)
			importAPI.GET("/status", importHandler.GetImportStatus)
		}
	}
}

// ginLogger 自定义日志中间件
func ginLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		logger.Info("HTTP请求",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
		)
	}
}

