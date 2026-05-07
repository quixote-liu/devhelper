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

$(BACKEND_DIR)/.env:
	cp $(BACKEND_DIR)/.env.example $(BACKEND_DIR)/.env

backend: $(BACKEND_DIR)/.env ## 启动后端（开发模式）
	cd $(BACKEND_DIR) && go run cmd/server/main.go

frontend: ## 启动前端（开发模式）
	cd $(FRONTEND_DIR) && npm run dev

# --- 构建 ---

build: build-backend build-frontend ## 构建前后端

build-backend: ## 构建后端二进制
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go
	@# 复制 .env 到 build/ 并调整路径配置
	cp $(BACKEND_DIR)/.env.example $(BUILD_DIR)/.env
	@# 修正 build/.env 中的相对路径（适配从 build/ 目录运行）
	@sed -i 's|DB_PATH=./devhelper.db|DB_PATH=./devhelper.db|g' $(BUILD_DIR)/.env
	@sed -i 's|STATIC_FILES_PATH=../frontend/dist|STATIC_FILES_PATH=./static|g' $(BUILD_DIR)/.env

build-frontend: ## 构建前端静态文件
	cp $(FRONTEND_DIR)/.env.example $(FRONTEND_DIR)/.env
	cd $(FRONTEND_DIR) && npm run build
	mkdir -p $(BUILD_DIR)/static
	cp -r $(FRONTEND_DIR)/dist/. $(BUILD_DIR)/static/

build-prod: build-frontend build-backend ## 生产构建（前端+后端，产物统一在 build/）

# --- 运行 ---

run-prod: build-prod ## 生产模式运行（后端托管前端）
	@echo "Starting in production mode..."
	cd $(BUILD_DIR) && SERVE_STATIC=true STATIC_FILES_PATH=./static ./$(BINARY_NAME)

# --- 依赖安装 ---

install: ## 安装前后端依赖
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && npm install

install-backend: ## 安装后端依赖
	cd $(BACKEND_DIR) && go mod tidy

install-frontend: ## 安装前端依赖
	cd $(FRONTEND_DIR) && npm install

# --   - 清理 ---

clean: ## 清理构建产物
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm $(BACKEND_DIR)/.env
	rm $(FRONTEND_DIR)/.env

# --- 帮助 ---

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
