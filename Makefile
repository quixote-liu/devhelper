.PHONY: dev backend frontend build build-backend build-frontend build-prod run-prod clean install lint test help

# 默认目标
.DEFAULT_GOAL := help

# 变量
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BUILD_DIR := build
BINARY_NAME := devhelper

# --- 开发 ---

dev: ## 同时启动前后端（开发模式）
	@echo "Starting backend and frontend..."
	$(MAKE) backend & $(MAKE) frontend & wait

backend: ## 启动后端（开发模式）
	cd $(BACKEND_DIR) && go run cmd/server/main.go

frontend: ## 启动前端（开发模式）
	cd $(FRONTEND_DIR) && npm run dev

# --- 构建 ---

build: build-backend build-frontend ## 构建前后端

build-backend: ## 构建后端二进制
	cd $(BACKEND_DIR) && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

build-frontend: ## 构建前端静态文件
	cd $(FRONTEND_DIR) && npm run build

build-prod: build-frontend build-backend ## 生产构建（前端+后端）

# --- 运行 ---

run-prod: build-prod ## 生产模式运行（后端托管前端）
	@echo "Starting in production mode..."
	cd $(BACKEND_DIR) && SERVE_STATIC=true STATIC_FILES_PATH=../$(FRONTEND_DIR)/dist ../$(BUILD_DIR)/$(BINARY_NAME)

# --- 依赖安装 ---

install: ## 安装前后端依赖
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && npm install

install-backend: ## 安装后端依赖
	cd $(BACKEND_DIR) && go mod download

install-frontend: ## 安装前端依赖
	cd $(FRONTEND_DIR) && npm install

# --- 代码质量 ---

lint: ## 运行前后端 lint
	cd $(BACKEND_DIR) && go vet ./...
	cd $(FRONTEND_DIR) && npm run lint

lint-backend: ## 运行后端 lint
	cd $(BACKEND_DIR) && go vet ./...

lint-frontend: ## 运行前端 lint
	cd $(FRONTEND_DIR) && npm run lint

test: ## 运行后端测试
	cd $(BACKEND_DIR) && go test ./...

# --   - 清理 ---

clean: ## 清理构建产物
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist

# --- 帮助 ---

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
