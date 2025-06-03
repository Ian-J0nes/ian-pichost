# Makefile for Ian·PicHost Image Host Service

# 变量定义
APP_NAME := ian-pichost
BINARY_NAME := ian-pichost
DOCKER_IMAGE := ian-pichost
DOCKER_TAG := latest
GO_VERSION := 1.24

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "Ian·PicHost Image Host Service"
	@echo "=========================="
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 开发相关
.PHONY: dev
dev: ## 启动开发服务器
	@echo "启动开发服务器..."
	go run main.go

.PHONY: build
build: ## 编译二进制文件
	@echo "编译二进制文件..."
	go build -o $(BINARY_NAME) main.go

.PHONY: build-linux
build-linux: ## 编译 Linux 二进制文件
	@echo "编译 Linux 二进制文件..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux main.go

.PHONY: build-windows
build-windows: ## 编译 Windows 二进制文件
	@echo "编译 Windows 二进制文件..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe main.go

.PHONY: build-all
build-all: build-linux build-windows ## 编译所有平台的二进制文件
	@echo "所有平台编译完成"

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download

.PHONY: tidy
tidy: ## 整理依赖
	@echo "整理依赖..."
	go mod tidy

.PHONY: vendor
vendor: ## 创建 vendor 目录
	@echo "创建 vendor 目录..."
	go mod vendor

# 测试相关
.PHONY: test
test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码质量
.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

.PHONY: vet
vet: ## 代码静态检查
	@echo "代码静态检查..."
	go vet ./...

.PHONY: lint
lint: ## 代码 lint 检查 (需要安装 golangci-lint)
	@echo "代码 lint 检查..."
	golangci-lint run

# Docker 相关
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run -d -p 8080:8080 -v $(PWD)/uploads:/app/uploads --name $(APP_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-stop
docker-stop: ## 停止 Docker 容器
	@echo "停止 Docker 容器..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

.PHONY: docker-logs
docker-logs: ## 查看 Docker 容器日志
	@echo "查看 Docker 容器日志..."
	docker logs -f $(APP_NAME)

.PHONY: docker-shell
docker-shell: ## 进入 Docker 容器
	@echo "进入 Docker 容器..."
	docker exec -it $(APP_NAME) /bin/sh

# Docker Compose 相关
.PHONY: up
up: ## 启动 docker-compose 服务
	@echo "启动 docker-compose 服务..."
	docker-compose up -d

.PHONY: down
down: ## 停止 docker-compose 服务
	@echo "停止 docker-compose 服务..."
	docker-compose down

.PHONY: logs
logs: ## 查看 docker-compose 日志
	@echo "查看 docker-compose 日志..."
	docker-compose logs -f

.PHONY: restart
restart: down up ## 重启 docker-compose 服务

# 清理相关
.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME).exe
	rm -f coverage.out
	rm -f coverage.html

.PHONY: clean-docker
clean-docker: ## 清理 Docker 镜像和容器
	@echo "清理 Docker 镜像和容器..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

.PHONY: clean-all
clean-all: clean clean-docker ## 清理所有文件

# 部署相关
.PHONY: deploy
deploy: docker-build docker-stop docker-run ## 部署应用 (构建镜像并运行)

# 初始化项目
.PHONY: init
init: deps tidy ## 初始化项目
	@echo "项目初始化完成"

# 检查环境
.PHONY: check
check: ## 检查开发环境
	@echo "检查 Go 版本..."
	@go version
	@echo "检查 Docker 版本..."
	@docker --version
	@echo "检查 Docker Compose 版本..."
	@docker-compose --version

# 生产环境构建
.PHONY: release
release: clean fmt vet test build-all docker-build ## 发布版本 (完整构建流程)
	@echo "发布版本完成"
