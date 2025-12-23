package main

import (
	"log"
	"os"
	"time"

	"gitlab-webhook-server/internal/config"
	"gitlab-webhook-server/internal/database"
	"gitlab-webhook-server/internal/gitlab"
	"gitlab-webhook-server/internal/handler"
	"gitlab-webhook-server/internal/logger"
	"gitlab-webhook-server/internal/middleware"
	"gitlab-webhook-server/internal/queue"
	"gitlab-webhook-server/internal/router"
	"gitlab-webhook-server/internal/service/commit"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	// 初始化数据库
	if err := database.Init(cfg, zapLogger); err != nil {
		zapLogger.Fatal("数据库初始化失败", zap.Error(err))
	}
	defer database.Close()

	// 执行数据库迁移
	if err := database.Migrate(); err != nil {
		zapLogger.Fatal("数据库迁移失败", zap.Error(err))
	}

	// 设置 Gin 模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建工作池
	workerPool := queue.NewWorkerPool(
		cfg.WorkerPool.Workers,
		cfg.WorkerPool.QueueSize,
		zapLogger,
	)
	workerPool.Start()
	defer workerPool.Stop()

	// 创建限流器
	rateLimitWindow, err := time.ParseDuration(cfg.RateLimit.Window)
	if err != nil {
		rateLimitWindow = time.Minute
		zapLogger.Warn("解析限流窗口失败，使用默认值 1m", zap.Error(err))
	}
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Limit, rateLimitWindow, zapLogger)

	// 创建路由
	r := router.New(zapLogger)

	// 应用限流中间件
	r.Use(rateLimiter.Limit())

	// 创建提交服务
	commitService := commit.NewCommitServiceV2(database.DB, zapLogger)

	// 创建 GitLab 客户端（如果配置了）
	var gitlabClient *gitlab.Client
	var importHandler *handler.ImportHandler
	if cfg.GitLab.BaseURL != "" && cfg.GitLab.Token != "" {
		client, err := gitlab.NewClient(cfg.GitLab.BaseURL, cfg.GitLab.Token, zapLogger)
		if err != nil {
			zapLogger.Warn("GitLab 客户端初始化失败，历史数据导入功能将不可用", zap.Error(err))
		} else {
			gitlabClient = client
			importHandler = handler.NewImportHandler(gitlabClient, commitService, database.DB, zapLogger)
			zapLogger.Info("GitLab 客户端初始化成功")
		}
	}

	// 注册路由
	webhookHandler := handler.NewWebhookHandler(database.DB, workerPool, zapLogger)
	statsHandler := handler.NewStatsHandler(database.DB, zapLogger)
	router.RegisterRoutes(r, webhookHandler, statsHandler, importHandler)

	// 启动服务器
	addr := ":" + cfg.Port
	zapLogger.Infof("🚀 服务器启动在端口 %s", cfg.Port)
	zapLogger.Infof("📡 Webhook 端点: http://localhost%s/webhook", addr)
	zapLogger.Infof("💚 健康检查: http://localhost%s/health", addr)

	if err := r.Run(addr); err != nil {
		zapLogger.Fatalf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}

