package handler

import (
	"net/http"
	"time"

	"gitlab-webhook-server/internal/gitlab"
	"gitlab-webhook-server/internal/service"
	"gitlab-webhook-server/internal/service/commit"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ImportHandler 历史数据导入处理器
type ImportHandler struct {
	logger        *zap.Logger
	importService *service.ImportService
}

// NewImportHandler 创建新的导入处理器
func NewImportHandler(
	gitlabClient *gitlab.Client,
	commitService *commit.CommitServiceV2,
	db *gorm.DB,
	logger *zap.Logger,
) *ImportHandler {
	importService := service.NewImportService(gitlabClient, commitService, db, logger)
	return &ImportHandler{
		logger:        logger,
		importService: importService,
	}
}

// ImportProject 导入项目的提交记录
// POST /api/import/project
// Body: {"project_id": "123", "since": "2024-01-01T00:00:00Z", "until": "2024-12-31T23:59:59Z", "batch_size": 100}
func (h *ImportHandler) ImportProject(c *gin.Context) {
	var req struct {
		ProjectID string `json:"project_id" binding:"required"`
		Since     string `json:"since"`
		Until     string `json:"until"`
		BatchSize int    `json:"batch_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("解析请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 解析时间
	var since, until *time.Time
	if req.Since != "" {
		if t, err := time.Parse(time.RFC3339, req.Since); err == nil {
			since = &t
		}
	}
	if req.Until != "" {
		if t, err := time.Parse(time.RFC3339, req.Until); err == nil {
			until = &t
		}
	}

	// 设置默认批次大小
	batchSize := req.BatchSize
	if batchSize == 0 {
		batchSize = 100
	}

	// 异步导入
	go func() {
		result, err := h.importService.ImportProjectCommits(
			req.ProjectID,
			since,
			until,
			batchSize,
		)
		if err != nil {
			h.logger.Error("导入失败", zap.Error(err))
			return
		}

		h.logger.Info("导入完成",
			zap.String("project_id", result.ProjectID),
			zap.Int("imported", result.Imported),
			zap.Int("failed", result.Failed),
		)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "导入任务已启动",
		"project_id": req.ProjectID,
		"status":     "processing",
	})
}

// GetImportStatus 获取导入状态（简化版，实际可以维护导入任务状态）
// GET /api/import/status?project_id=123
func (h *ImportHandler) GetImportStatus(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	// TODO: 实现导入状态查询
	// 可以维护一个导入任务状态表，或者通过查询数据库中的提交记录来判断

	c.JSON(http.StatusOK, gin.H{
		"project_id": projectID,
		"status":     "completed", // 简化实现
		"message":    "查询导入状态功能待实现",
	})
}

