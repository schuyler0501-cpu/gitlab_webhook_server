# AI 辅助开发指南

## 🎯 本指南的目的

这份文档帮助你更好地与 AI 协作，提高开发效率和代码质量。

## 📖 项目结构说明

### 核心目录

- **`cmd/server/`**: 应用入口，包含 `main.go`
- **`internal/handler/`**: HTTP 请求处理器，处理路由和请求验证
- **`internal/service/`**: 业务逻辑实现，核心功能代码
- **`internal/model/`**: 数据模型定义，结构体和类型
- **`internal/router/`**: 路由定义和中间件配置
- **`internal/config/`**: 配置管理，环境变量和配置加载
- **`internal/logger/`**: 日志工具，统一的日志输出

### 为什么这样组织？

1. **清晰的分层**: Handler → Service → Repository，职责明确
2. **易于测试**: 每层都可以独立测试
3. **便于 AI 理解**: 标准化的 Go 项目结构让 AI 更容易找到相关代码
4. **符合 Go 惯例**: 遵循 Go 社区的最佳实践

## 💡 与 AI 协作的最佳实践

### 1. 提供清晰的上下文

**✅ 好的方式**:
```
在 internal/service/commit/commit_service.go 的 RecordCommit 方法中，
添加数据库持久化逻辑，使用 GORM 和 PostgreSQL 存储提交记录。
```

**❌ 不好的方式**:
```
帮我保存数据
```

### 2. 引用现有代码

**✅ 好的方式**:
```
参考 internal/service/webhook_service.go 中的 handlePushEvent 方法，
在 commit_service.go 中添加类似的错误处理逻辑。
```

**❌ 不好的方式**:
```
看看之前的代码，照着写
```

### 3. 分步骤实现

**✅ 好的方式**:
```
第一步：在 internal/model/ 中定义数据库模型
第二步：创建 internal/repository/commit_repository.go
第三步：在 commit_service.go 中集成 repository
```

**❌ 不好的方式**:
```
实现完整的数据库功能
```

### 4. 明确技术栈和依赖

**✅ 好的方式**:
```
使用 GORM 和 PostgreSQL，参考项目现有的配置风格。
```

**❌ 不好的方式**:
```
用数据库保存
```

## 🔧 常用开发场景

### 场景 1: 添加新的 API 端点

**提示词示例**:
```
在 internal/router/router.go 中添加一个新的 GET 端点 /api/stats，
返回最近 7 天的提交统计信息。
需要：
1. 在 router 中定义路由
2. 在 handler 中创建处理方法
3. 在 service 中实现统计逻辑
```

### 场景 2: 添加新的 Webhook 事件处理

**提示词示例**:
```
在 internal/service/webhook_service.go 中添加对 'Merge Request Hook' 事件的处理。
参考现有的 handlePushEvent 方法的结构。
```

### 场景 3: 重构代码

**提示词示例**:
```
重构 internal/service/commit/commit_service.go 中的 RecordCommit 方法，
将数据持久化逻辑提取到独立的 repository 层。
保持现有功能不变，只改变代码结构。
```

### 场景 4: 添加数据库支持

**提示词示例**:
```
为项目添加 PostgreSQL 数据库支持：
1. 在 go.mod 中添加 GORM 依赖
2. 在 internal/config/ 中添加数据库配置
3. 创建 internal/repository/ 目录和数据库连接
4. 在 commit_service 中集成数据库操作
```

## 📝 代码质量检查清单

在提交代码前，确保：

- [ ] 运行 `make lint` 没有错误
- [ ] 运行 `make fmt` 格式化代码
- [ ] 运行 `go test ./...` 所有测试通过
- [ ] 所有新功能都有相应的测试
- [ ] 代码遵循 Go 的命名规范
- [ ] 错误处理完善，没有忽略错误

## 🚀 快速开始示例

### 添加新功能的标准流程

1. **定义模型** (在 `internal/model/`)
```go
type NewFeature struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

2. **实现服务** (在 `internal/service/`)
```go
type NewFeatureService struct {
    logger *zap.Logger
}

func (s *NewFeatureService) DoSomething() error {
    // 实现逻辑
    return nil
}
```

3. **创建处理器** (在 `internal/handler/`)
```go
func (h *Handler) HandleRequest(c *gin.Context) {
    // 处理请求
}
```

4. **添加路由** (在 `internal/router/`)
```go
r.GET("/new-feature", handler.HandleRequest)
```

## 🎓 学习资源

- **Go 语言**: 理解 Go 的并发模型和错误处理
- **Gin 框架**: 理解中间件和路由机制
- **GitLab Webhooks**: 了解不同事件类型和数据结构
- **Zap 日志**: 高效的日志记录

## 💬 与 AI 对话的技巧

1. **一次一个任务**: 不要在一个提示中包含太多需求
2. **提供反馈**: 如果 AI 的实现不符合预期，明确指出问题
3. **迭代改进**: 先实现基础功能，再逐步完善
4. **保持上下文**: 在对话中引用之前的代码和决策
5. **使用 Go 术语**: 使用 "struct", "interface", "goroutine" 等 Go 特定术语

## 🔍 Go 特定提示

- **接口设计**: Go 的接口是隐式实现的，设计小而专注的接口
- **错误处理**: 始终检查错误，使用 `fmt.Errorf` 添加上下文
- **并发安全**: 如果使用 goroutine，注意并发安全
- **包设计**: 保持包的职责单一，避免循环依赖

---

记住：**清晰的沟通 = 更好的代码** 🎯
