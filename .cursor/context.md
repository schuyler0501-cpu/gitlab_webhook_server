# 项目上下文文档

这个文档为 Cursor AI 提供项目的整体上下文信息，帮助 AI 更好地理解项目结构和业务逻辑。

## 项目背景

这是一个 GitLab Webhook 服务器，主要用于：

1. **接收 GitLab Webhook 事件**: 监听代码推送、标签推送等事件
2. **记录代码提交**: 解析提交信息，记录提交者、时间、文件变更等
3. **效能度量**: 统计团队成员的代码产出，用于评估工作饱和度和任务完成情况

## 核心业务流程

### Webhook 处理流程

```
GitLab 推送代码
    ↓
触发 Webhook 事件
    ↓
服务器接收请求 (handler/webhook_handler.go)
    ↓
验证 Token 和解析请求
    ↓
调用服务层处理 (service/webhook_service.go)
    ↓
解析提交信息
    ↓
调用提交服务记录 (service/commit/commit_service.go)
    ↓
保存到数据存储（待实现）
    ↓
返回成功响应
```

### 数据模型

**CommitRecord** (model/commit.go):
- CommitID: 提交 ID
- Message: 提交信息
- Timestamp: 提交时间
- Author: 作者姓名
- AuthorEmail: 作者邮箱
- URL: 提交链接
- ProjectName: 项目名称
- ProjectPath: 项目路径
- AddedFiles: 新增文件列表
- ModifiedFiles: 修改文件列表
- RemovedFiles: 删除文件列表

## 技术决策

### 为什么选择 Gin？

- 轻量级，性能好
- 中间件支持完善
- 社区活跃，文档完善

### 为什么选择 Zap？

- 高性能的结构化日志
- 丰富的日志级别
- 适合生产环境

### 为什么使用 internal/ 目录？

- 防止外部包导入内部实现
- 明确包的可见性
- 符合 Go 项目最佳实践

## 待实现功能

1. **数据持久化**: 
   - 当前提交记录只打印日志，需要实现数据库存储
   - 建议使用 PostgreSQL + GORM

2. **统计查询 API**:
   - 按成员查询提交记录
   - 按时间范围统计
   - 按项目统计

3. **更多 Webhook 事件**:
   - Merge Request Hook
   - Issue Hook
   - Pipeline Hook

4. **认证和授权**:
   - API 密钥认证
   - 权限控制

## 常见问题

### Q: 如何添加新的 Webhook 事件处理？

A: 在 `internal/service/webhook_service.go` 的 `ProcessWebhook` 方法中添加新的 case，然后创建对应的处理方法。

### Q: 如何添加新的 API 端点？

A: 
1. 在 `internal/router/router.go` 的 `RegisterRoutes` 中添加路由
2. 在 `internal/handler/` 创建对应的处理器
3. 在 `internal/service/` 实现业务逻辑

### Q: 如何添加数据库支持？

A:
1. 在 `go.mod` 添加数据库驱动和 ORM（如 GORM）
2. 在 `internal/config/` 添加数据库配置
3. 创建 `internal/repository/` 目录
4. 在 service 层集成 repository

## 相关文件

- **路由定义**: `internal/router/router.go`
- **Webhook 处理**: `internal/handler/webhook_handler.go`
- **业务逻辑**: `internal/service/webhook_service.go`
- **提交服务**: `internal/service/commit/commit_service.go`
- **数据模型**: `internal/model/commit.go`
- **配置管理**: `internal/config/config.go`

## 开发工作流

1. 修改代码
2. 运行 `make fmt` 格式化
3. 运行 `make lint` 检查
4. 运行 `make test` 测试
5. 提交代码

## 性能考虑

- Webhook 处理应该是异步的，避免阻塞
- 日志记录使用异步方式
- 数据库操作使用连接池
- 考虑使用消息队列处理高并发场景

