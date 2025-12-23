# 常用提示词模板

这个文件包含一些常用的提示词模板，帮助用户更好地与 AI 协作。

## 功能开发

### 添加新的 API 端点

```
在项目中添加一个新的 API 端点 /api/stats，用于返回最近 7 天的提交统计信息。

要求：
1. 在 internal/router/router.go 中添加路由
2. 在 internal/handler/ 创建 stats_handler.go
3. 在 internal/service/ 创建 stats_service.go
4. 返回 JSON 格式的统计数据，包括：
   - 总提交数
   - 按作者分组的提交数
   - 按项目分组的提交数
5. 添加适当的错误处理和日志记录
6. 遵循项目的代码规范和架构模式
```

### 添加数据库支持

```
为项目添加 PostgreSQL 数据库支持，用于持久化提交记录。

要求：
1. 在 go.mod 中添加 GORM 和 PostgreSQL 驱动依赖
2. 在 internal/config/ 中添加数据库配置（连接字符串、连接池等）
3. 创建 internal/repository/ 目录和 commit_repository.go
4. 在 internal/service/commit/commit_service.go 中集成 repository
5. 实现数据库迁移逻辑
6. 添加数据库连接的健康检查
7. 遵循项目的错误处理和日志规范
```

### 添加新的 Webhook 事件

```
在项目中添加对 GitLab 'Merge Request Hook' 事件的处理。

要求：
1. 在 internal/service/webhook_service.go 的 ProcessWebhook 中添加新的 case
2. 创建 handleMergeRequestEvent 方法
3. 解析 Merge Request 的相关信息（作者、状态、变更等）
4. 记录 Merge Request 信息（可以创建新的 model）
5. 添加适当的日志记录
6. 参考现有的 handlePushEvent 方法的结构
```

## 代码重构

### 提取公共逻辑

```
重构 internal/service/webhook_service.go 中的代码，将 JSON 解析逻辑提取到独立的工具函数中。

要求：
1. 创建 internal/utils/json_parser.go
2. 提取 parseCommit 和 parseStringSlice 方法
3. 保持现有功能不变
4. 添加单元测试
5. 更新相关注释
```

### 优化错误处理

```
优化项目中的错误处理，统一错误响应格式。

要求：
1. 创建 internal/errors/errors.go，定义统一的错误类型
2. 在 handler 层使用统一的错误处理中间件
3. 所有错误响应使用统一的格式
4. 保持向后兼容
```

## 测试

### 添加单元测试

```
为 internal/service/commit/commit_service.go 添加完整的单元测试。

要求：
1. 测试 RecordCommit 方法的所有分支
2. 测试 GetMemberCommits 方法
3. 使用表驱动测试
4. 使用 mock 对象隔离依赖
5. 测试覆盖率要达到 80%+
6. 测试正常情况和异常情况
```

### 添加集成测试

```
为 Webhook 处理流程添加集成测试。

要求：
1. 创建 test/integration/webhook_test.go
2. 测试完整的 Webhook 处理流程
3. 使用测试 HTTP 服务器
4. 验证请求和响应
5. 测试各种错误场景
```

## 性能优化

### 优化 Webhook 处理性能

```
优化 Webhook 处理性能，支持高并发场景。

要求：
1. 使用 goroutine 异步处理 Webhook
2. 使用 channel 进行任务队列管理
3. 添加限流机制
4. 优化日志记录（使用异步日志）
5. 添加性能监控指标
6. 保持代码清晰和可维护
```

## 代码审查

### 审查代码质量

```
请审查以下代码，检查是否符合项目规范：

[粘贴代码]

检查项：
1. 命名规范
2. 错误处理
3. 日志记录
4. 代码注释
5. 架构模式
6. 性能考虑
```

## 调试

### 调试 Webhook 问题

```
我在处理 GitLab Webhook 时遇到问题：[描述问题]

相关代码：
- internal/handler/webhook_handler.go
- internal/service/webhook_service.go

请帮我：
1. 分析可能的原因
2. 提供调试建议
3. 修复代码（如果发现问题）
```

## 文档

### 更新 API 文档

```
为项目添加 API 文档，使用 Swagger/OpenAPI。

要求：
1. 添加 swagger 注释到所有 API 端点
2. 生成 OpenAPI 规范文件
3. 添加 Swagger UI 端点
4. 文档要包含请求/响应示例
```

---

**提示**: 使用这些模板时，根据实际情况调整具体需求。越详细的描述，AI 生成的代码质量越高。

