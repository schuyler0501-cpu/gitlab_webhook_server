.PHONY: build run dev test lint fmt clean help

# 变量定义
APP_NAME=gitlab-webhook-server
BIN_DIR=bin
MAIN_PATH=cmd/server/main.go

# 构建应用
build:
	@echo "构建应用..."
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BIN_DIR)/$(APP_NAME)"

# 运行应用
run: build
	@echo "启动应用..."
	@./$(BIN_DIR)/$(APP_NAME)

# 开发模式（使用 air 热重载，如果安装了的话）
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air 未安装，使用普通模式运行..."; \
		go run $(MAIN_PATH); \
	fi

# 运行测试
test:
	@echo "运行测试..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "测试完成，覆盖率报告: coverage.html"

# 代码检查
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过..."; \
		echo "安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@goimports -w .
	@echo "格式化完成"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "清理完成"

# 安装依赖
deps:
	@echo "下载依赖..."
	@go mod download
	@go mod tidy
	@echo "依赖安装完成"

# 显示帮助信息
help:
	@echo "可用命令:"
	@echo "  make build    - 构建应用"
	@echo "  make run      - 构建并运行应用"
	@echo "  make dev      - 开发模式（热重载）"
	@echo "  make test     - 运行测试"
	@echo "  make lint     - 代码检查"
	@echo "  make fmt      - 格式化代码"
	@echo "  make clean    - 清理构建文件"
	@echo "  make deps     - 安装依赖"
	@echo "  make help     - 显示帮助信息"

