package handler

import (
	"net/http"
	"time"

	"gitlab-webhook-server/internal/service/commit"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	logger        *zap.Logger
	commitService *commit.CommitServiceV2
}

// NewStatsHandler 创建新的统计处理器
func NewStatsHandler(db *gorm.DB, logger *zap.Logger) *StatsHandler {
	return &StatsHandler{
		logger:        logger,
		commitService: commit.NewCommitServiceV2(db, logger),
	}
}

// GetMemberStats 获取成员统计信息
// GET /api/stats/member?email=user@example.com&start_date=2024-01-01&end_date=2024-02-01
func (h *StatsHandler) GetMemberStats(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email 参数必填"})
		return
	}

	// 解析时间参数
	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			// 设置为当天的结束时间
			t = t.Add(24*time.Hour - time.Second)
			endDate = &t
		}
	}

	stats, err := h.commitService.GetMemberStats(email, startDate, endDate)
	if err != nil {
		h.logger.Error("获取成员统计失败",
			zap.Error(err),
			zap.String("email", email),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":        email,
		"commit_count": stats.CommitCount,
		"total_added":  stats.TotalAdded,
		"total_removed": stats.TotalRemoved,
		"total_files":  stats.TotalFiles,
	})
}

// GetLanguageStats 获取语言统计信息
// GET /api/stats/languages?email=user@example.com&start_date=2024-01-01&end_date=2024-02-01
func (h *StatsHandler) GetLanguageStats(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email 参数必填"})
		return
	}

	// 解析时间参数
	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endDate = &t
		}
	}

	stats, err := h.commitService.GetLanguageStats(email, startDate, endDate)
	if err != nil {
		h.logger.Error("获取语言统计失败",
			zap.Error(err),
			zap.String("email", email),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":    email,
		"languages": stats,
	})
}

// GetMemberCommits 获取成员提交记录
// GET /api/stats/commits?email=user@example.com&start_date=2024-01-01&end_date=2024-02-01
func (h *StatsHandler) GetMemberCommits(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email 参数必填"})
		return
	}

	// 解析时间参数
	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endDate = &t
		}
	}

	commits, err := h.commitService.GetMemberCommits(email, startDate, endDate)
	if err != nil {
		h.logger.Error("获取成员提交记录失败",
			zap.Error(err),
			zap.String("email", email),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取提交记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":   email,
		"commits": commits,
		"count":   len(commits),
	})
}

