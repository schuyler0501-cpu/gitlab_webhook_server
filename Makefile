.PHONY: build run dev test lint fmt clean help deps install-tools check-env version

# 变量定义
APP_NAME=gitlab-webhook-server
BIN_DIR=bin
MAIN_PATH=cmd/server/main.go
GO_VERSION=1.21

# 检测操作系统
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RM := powershell -Command "if (Test-Path '$(1)') { Remove-Item -Recurse -Force '$(1)' }"
    MKDIR := powershell -Command "New-Item -ItemType Directory -Force -Path '$(1)' | Out-Null"
    CHECK_CMD := where
    PATH_SEP := ;
else
    DETECTED_OS := $(shell uname -s)
    RM := rm -rf
    MKDIR := mkdir -p
    CHECK_CMD := command -v
    PATH_SEP := :
endif

# 检查命令是否存在（跨平台）
define check_command
	@$(CHECK_CMD) $(1) > /dev/null 2>&1 || (echo "错误: $(1) 未安装" && exit 1)
endef

# 构建应用
build:
	@echo "🔨 构建应用..."
	@$(MKDIR) $(BIN_DIR)
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME)$(if $(filter Windows,$(DETECTED_OS)),.exe,) $(MAIN_PATH)
	@echo "✅ 构建完成: $(BIN_DIR)/$(APP_NAME)$(if $(filter Windows,$(DETECTED_OS)),.exe,)"

# 运行应用
run: build
	@echo "🚀 启动应用..."
	@$(if $(filter Windows,$(DETECTED_OS)),$(BIN_DIR)/$(APP_NAME).exe,$(BIN_DIR)/$(APP_NAME))

# 开发模式（使用 air 热重载，如果安装了的话）
dev:
	@echo "💻 启动开发模式..."
	@if $(CHECK_CMD) air > /dev/null 2>&1; then \
		echo "✅ 使用 Air 热重载..."; \
		air; \
	else \
		echo "⚠️  Air 未安装，使用普通模式运行..."; \
		echo "💡 安装 Air: make install-tools"; \
		go run $(MAIN_PATH); \
	fi

# 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v -race -coverprofile=coverage.out ./...
	@if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "✅ 测试完成，覆盖率报告: coverage.html"; \
	else \
		echo "⚠️  未生成覆盖率报告"; \
	fi

# 代码检查
lint:
	@echo "🔍 运行代码检查..."
	@if $(CHECK_CMD) golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
		echo "✅ 代码检查完成"; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过..."; \
		echo "💡 安装命令: make install-tools"; \
	fi

# 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	@go fmt ./...
	@if $(CHECK_CMD) goimports > /dev/null 2>&1; then \
		goimports -w .; \
		echo "✅ 格式化完成"; \
	else \
		echo "⚠️  goimports 未安装，跳过导入整理..."; \
		echo "💡 安装命令: make install-tools"; \
	fi

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@if [ -d $(BIN_DIR) ]; then $(RM) $(BIN_DIR); fi
	@if [ -f coverage.out ]; then $(RM) coverage.out; fi
	@if [ -f coverage.html ]; then $(RM) coverage.html; fi
	@go clean
	@echo "✅ 清理完成"

# 安装依赖
deps:
	@echo "📦 下载依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖安装完成"

# 安装开发工具
install-tools:
	@echo "🛠️  安装开发工具..."
	@echo "安装 Air (热重载)..."
	@go install github.com/air-verse/air@latest
	@echo "安装 golangci-lint (代码检查)..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "安装 goimports (导入管理)..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ 工具安装完成"
	@echo "💡 请确保 $(shell go env GOPATH)/bin 在 PATH 环境变量中"

# 检查开发环境
check-env:
	@echo "🔍 检查开发环境..."
	@echo "操作系统: $(DETECTED_OS)"
	@echo "Go 版本: $(shell go version 2>/dev/null || echo '未安装')"
	@echo "Go 路径: $(shell go env GOPATH 2>/dev/null || echo '未设置')"
	@echo ""
	@echo "工具检查:"
	@echo -n "  Air: "
	@if $(CHECK_CMD) air > /dev/null 2>&1; then echo "✅ 已安装"; else echo "❌ 未安装 (运行: make install-tools)"; fi
	@echo -n "  golangci-lint: "
	@if $(CHECK_CMD) golangci-lint > /dev/null 2>&1; then echo "✅ 已安装"; else echo "❌ 未安装 (运行: make install-tools)"; fi
	@echo -n "  goimports: "
	@if $(CHECK_CMD) goimports > /dev/null 2>&1; then echo "✅ 已安装"; else echo "❌ 未安装 (运行: make install-tools)"; fi
	@echo ""
	@echo "项目检查:"
	@echo -n "  .env 文件: "
	@if [ -f .env ]; then echo "✅ 存在"; else echo "❌ 不存在 (从 env.example 复制)"; fi
	@echo -n "  go.mod: "
	@if [ -f go.mod ]; then echo "✅ 存在"; else echo "❌ 不存在"; fi

# 显示版本信息
version:
	@echo "📋 项目信息:"
	@echo "  项目名称: $(APP_NAME)"
	@echo "  Go 版本: $(shell go version 2>/dev/null || echo '未知')"
	@echo "  操作系统: $(DETECTED_OS)"
	@if [ -f go.mod ]; then \
		echo "  模块路径: $(shell grep '^module' go.mod | awk '{print $$2}')"; \
		echo "  Go 版本要求: >= $(GO_VERSION)"; \
	fi

# 显示帮助信息
help:
	@echo "📚 可用命令:"
	@echo ""
	@echo "  🏗️  构建和运行:"
	@echo "    make build        - 构建应用"
	@echo "    make run          - 构建并运行应用"
	@echo "    make dev          - 开发模式（热重载，需要 Air）"
	@echo ""
	@echo "  🧪 测试和检查:"
	@echo "    make test         - 运行测试并生成覆盖率报告"
	@echo "    make lint         - 代码检查（需要 golangci-lint）"
	@echo "    make fmt          - 格式化代码（需要 goimports）"
	@echo ""
	@echo "  🛠️  工具和依赖:"
	@echo "    make deps         - 安装/更新依赖"
	@echo "    make install-tools - 安装开发工具（Air, golangci-lint, goimports）"
	@echo "    make check-env    - 检查开发环境"
	@echo ""
	@echo "  🧹 清理:"
	@echo "    make clean        - 清理构建文件和缓存"
	@echo ""
	@echo "  ℹ️  信息:"
	@echo "    make version      - 显示版本信息"
	@echo "    make help         - 显示此帮助信息"
	@echo ""
	@echo "💡 提示:"
	@echo "  - 首次使用建议运行: make install-tools && make check-env"
	@echo "  - Windows 用户需要安装 Make 工具（如 Git Bash 或 Chocolatey）"
	@echo "  - 确保 Go 工具路径在 PATH 中: $(shell go env GOPATH 2>/dev/null || echo '$GOPATH')/bin"
