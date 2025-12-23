# Cursor AI 指令

这个文件包含给 Cursor AI 的特定指令，帮助 AI 更好地理解和协助开发。

## ⚠️ 重要：记忆系统

**在开始任何对话前，必须执行以下步骤**：

1. **读取 `.cursor/memory.md` 文件** - 了解用户的重要指令、偏好和项目特定要求
2. **读取 `.cursor/conversations.md` 文件** - 回顾历史有效对话，建立长久持续的记忆
3. **应用记忆内容** - 记忆文件中的内容优先于本文件中的通用指令

### 有效对话处理

当检测到用户消息包含 `[有效对话]` 标识时：

1. **立即识别**: 在回复前识别该标识
2. **获取北京时间**: 使用系统时间转换为北京时间（UTC+8）
3. **生成回复**: 正常生成回复内容
4. **保存对话**: 将对话保存到 `.cursor/conversations.md`，格式：
   ```markdown
   ## 对话记录 - [北京时间 YYYY-MM-DD HH:MM:SS]
   
   ### 用户
   [用户消息，去除 [有效对话] 标识]
   
   ### AI 回复
   [完整回复]
   
   ---
   ```
5. **确认保存**: 在回复末尾告知用户"✅ 此对话已保存到历史记录"

**注意**: 
- 保存时去除 `[有效对话]` 标识，只保存实际内容
- 时间格式必须使用北京时间
- 每次对话开始时都要回顾历史对话

## 代码生成规则

### 1. 始终遵循项目结构

- 新功能必须放在正确的目录层级
- Handler 只处理 HTTP 相关逻辑
- Service 包含业务逻辑
- Model 只定义数据结构

### 2. 错误处理

- 所有可能出错的操作都要检查错误
- 使用 `fmt.Errorf` 或 `errors.Wrap` 添加上下文
- 在 handler 层统一处理错误响应

### 3. 日志记录

- 使用结构化日志（zap）
- 包含足够的上下文信息
- 错误日志必须包含错误对象

### 4. 代码注释

- 所有公开的 API 必须有中文注释
- 复杂逻辑需要行内注释
- 注释要简洁明了

## 代码审查重点

当 AI 生成代码时，应该检查：

1. ✅ 是否遵循命名规范
2. ✅ 错误处理是否完善
3. ✅ 日志记录是否充分
4. ✅ 是否有适当的注释
5. ✅ 是否符合项目架构模式
6. ✅ 是否可以通过 lint 检查

## 常见任务模板

### 添加新的 API 端点

```go
// 1. 在 router.go 添加路由
r.GET("/api/endpoint", handler.HandleEndpoint)

// 2. 在 handler 创建方法
func (h *Handler) HandleEndpoint(c *gin.Context) {
    // 参数验证
    // 调用 service
    // 返回响应
}

// 3. 在 service 实现逻辑
func (s *Service) DoSomething() error {
    // 业务逻辑
}
```

### 添加新的 Webhook 事件

```go
// 在 webhook_service.go 的 ProcessWebhook 中添加
case "New Event":
    return s.handleNewEvent(payload)

// 创建处理方法
func (s *WebhookService) handleNewEvent(payload map[string]interface{}) error {
    // 处理逻辑
}
```

## AI 响应格式

当 AI 生成代码时，应该：

1. **解释代码**: 简要说明代码的作用
2. **提供上下文**: 说明代码在项目中的位置和作用
3. **提示下一步**: 告诉用户接下来需要做什么
4. **注意事项**: 提醒可能需要注意的点

## 代码风格示例

### ✅ 好的代码风格

```go
// HandleWebhook 处理 GitLab Webhook 请求
// 验证 token，解析请求体，并调用服务层处理业务逻辑
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
    eventType := c.GetHeader("X-Gitlab-Event")
    if eventType == "" {
        h.logger.Warn("缺少事件类型")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing event type"})
        return
    }

    var payload map[string]interface{}
    if err := c.ShouldBindJSON(&payload); err != nil {
        h.logger.Error("解析请求失败", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    if err := h.webhookService.ProcessWebhook(eventType, payload); err != nil {
        h.logger.Error("处理 Webhook 失败", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
```

### ❌ 不好的代码风格

```go
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
    // 缺少注释
    eventType := c.GetHeader("X-Gitlab-Event")
    var payload map[string]interface{}
    c.ShouldBindJSON(&payload)  // 忽略错误
    h.webhookService.ProcessWebhook(eventType, payload)  // 忽略错误
    c.JSON(200, gin.H{"ok": true})  // 使用魔法数字
}
```

## 特殊场景处理

### 数据库操作

```go
// ✅ 好的方式
func (r *Repository) SaveCommit(commit *model.CommitRecord) error {
    if err := r.db.Create(commit).Error; err != nil {
        return fmt.Errorf("保存提交记录失败: %w", err)
    }
    return nil
}

// ❌ 不好的方式
func (r *Repository) SaveCommit(commit *model.CommitRecord) {
    r.db.Create(commit)  // 忽略错误
}
```

### 并发处理

```go
// ✅ 使用 WaitGroup 和错误通道
var wg sync.WaitGroup
errCh := make(chan error, len(commits))

for _, commit := range commits {
    wg.Add(1)
    go func(c *model.CommitRecord) {
        defer wg.Done()
        if err := s.commitService.RecordCommit(c); err != nil {
            errCh <- err
        }
    }(commit)
}

wg.Wait()
close(errCh)
```

## 测试代码生成

生成测试代码时应该：

1. 使用表驱动测试
2. 测试正常情况和异常情况
3. 使用 mock 对象隔离依赖
4. 测试覆盖率要达到 80%+

```go
func TestWebhookHandler_HandleWebhook(t *testing.T) {
    tests := []struct {
        name           string
        setup          func() (*WebhookHandler, *gin.Context)
        expectedStatus int
    }{
        {
            name: "valid request",
            setup: func() (*WebhookHandler, *gin.Context) {
                // 设置测试环境
            },
            expectedStatus: http.StatusOK,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

## 文档生成

当添加新功能时，AI 应该：

1. 更新相关的文档
2. 在代码注释中说明功能
3. 提供使用示例

---

**记住**: 代码质量比速度更重要。生成的代码应该清晰、可维护、符合项目规范。

