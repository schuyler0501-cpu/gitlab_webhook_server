package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestWebhookHandler_Test(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试用的 logger
	logger, _ := zap.NewDevelopment()

	// 创建 handler
	handler := NewWebhookHandler(logger)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/webhook/test", nil)
	w := httptest.NewRecorder()

	// 创建 Gin 上下文
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 执行处理函数
	handler.Test(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，得到 %d", http.StatusOK, w.Code)
	}
}

