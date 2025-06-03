#!/bin/bash

# Ian·PicHost 部署脚本
# 用于快速部署 Ian·PicHost 图片托管服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
APP_NAME="ian-pichost"
DOCKER_IMAGE="ian-pichost:latest"
CONTAINER_NAME="ian-pichost"
PORT="8080"
UPLOAD_DIR="./uploads"
LOG_DIR="./logs"

# 函数：打印彩色消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}"
    echo "=================================="
    echo "  Ian·PicHost 部署脚本"
    echo "=================================="
    echo -e "${NC}"
}

# 函数：检查依赖
check_dependencies() {
    print_message "检查依赖..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_warning "Docker Compose 未安装，将使用 docker 命令部署"
        USE_COMPOSE=false
    else
        USE_COMPOSE=true
    fi
    
    print_message "依赖检查完成"
}

# 函数：创建必要的目录
create_directories() {
    print_message "创建必要的目录..."
    
    mkdir -p "$UPLOAD_DIR"
    mkdir -p "$LOG_DIR"
    
    print_message "目录创建完成"
}

# 函数：停止现有容器
stop_existing() {
    print_message "停止现有容器..."
    
    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        docker stop "$CONTAINER_NAME"
        docker rm "$CONTAINER_NAME"
        print_message "已停止现有容器"
    else
        print_message "没有运行中的容器"
    fi
}

# 函数：构建 Docker 镜像
build_image() {
    print_message "构建 Docker 镜像..."
    
    if [ -f "Dockerfile" ]; then
        docker build -t "$DOCKER_IMAGE" .
        print_message "镜像构建完成"
    else
        print_error "Dockerfile 不存在"
        exit 1
    fi
}

# 函数：使用 Docker Compose 部署
deploy_with_compose() {
    print_message "使用 Docker Compose 部署..."
    
    if [ -f "docker-compose.yml" ]; then
        docker-compose down
        docker-compose up -d
        print_message "Docker Compose 部署完成"
    else
        print_error "docker-compose.yml 不存在"
        exit 1
    fi
}

# 函数：使用 Docker 命令部署
deploy_with_docker() {
    print_message "使用 Docker 命令部署..."
    
    docker run -d \
        --name "$CONTAINER_NAME" \
        --restart unless-stopped \
        -p "$PORT:8080" \
        -v "$(pwd)/$UPLOAD_DIR:/app/uploads" \
        -v "$(pwd)/$LOG_DIR:/app/logs" \
        -e TZ=Asia/Shanghai \
        "$DOCKER_IMAGE"
    
    print_message "Docker 部署完成"
}

# 函数：检查部署状态
check_deployment() {
    print_message "检查部署状态..."
    
    sleep 5
    
    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        print_message "容器运行正常"
        
        # 检查健康状态
        if curl -f http://localhost:$PORT/ping -H "X-Auth-Key: your-key" &> /dev/null; then
            print_message "服务健康检查通过"
        else
            print_warning "服务健康检查失败，请检查配置"
        fi
    else
        print_error "容器启动失败"
        exit 1
    fi
}

# 函数：显示部署信息
show_deployment_info() {
    echo -e "${GREEN}"
    echo "=================================="
    echo "  部署完成！"
    echo "=================================="
    echo -e "${NC}"
    echo "服务地址: http://localhost:$PORT"
    echo "上传目录: $UPLOAD_DIR"
    echo "日志目录: $LOG_DIR"
    echo ""
    echo "常用命令:"
    echo "  查看日志: docker logs -f $CONTAINER_NAME"
    echo "  停止服务: docker stop $CONTAINER_NAME"
    echo "  重启服务: docker restart $CONTAINER_NAME"
    echo ""
    echo "配置说明:"
    echo "  请在 main.go 中修改 AuthKey 等配置"
    echo "  默认访问密钥: your-key"
}

# 函数：显示帮助信息
show_help() {
    echo "Ian·PicHost 部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -b, --build    构建镜像"
    echo "  -d, --deploy   部署服务"
    echo "  -s, --stop     停止服务"
    echo "  -r, --restart  重启服务"
    echo "  -l, --logs     查看日志"
    echo "  --clean        清理所有容器和镜像"
    echo ""
    echo "示例:"
    echo "  $0 -d          # 部署服务"
    echo "  $0 -b -d       # 构建并部署"
    echo "  $0 -r          # 重启服务"
}

# 函数：查看日志
show_logs() {
    print_message "查看服务日志..."
    docker logs -f "$CONTAINER_NAME"
}

# 函数：重启服务
restart_service() {
    print_message "重启服务..."
    docker restart "$CONTAINER_NAME"
    check_deployment
    print_message "服务重启完成"
}

# 函数：清理
clean_all() {
    print_warning "这将删除所有相关的容器和镜像，确定要继续吗？(y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_message "清理容器和镜像..."
        docker stop "$CONTAINER_NAME" 2>/dev/null || true
        docker rm "$CONTAINER_NAME" 2>/dev/null || true
        docker rmi "$DOCKER_IMAGE" 2>/dev/null || true
        print_message "清理完成"
    else
        print_message "取消清理"
    fi
}

# 主函数
main() {
    print_header
    
    # 解析命令行参数
    BUILD=false
    DEPLOY=false
    STOP=false
    RESTART=false
    LOGS=false
    CLEAN=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -b|--build)
                BUILD=true
                shift
                ;;
            -d|--deploy)
                DEPLOY=true
                shift
                ;;
            -s|--stop)
                STOP=true
                shift
                ;;
            -r|--restart)
                RESTART=true
                shift
                ;;
            -l|--logs)
                LOGS=true
                shift
                ;;
            --clean)
                CLEAN=true
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定任何选项，默认部署
    if [ "$BUILD" = false ] && [ "$DEPLOY" = false ] && [ "$STOP" = false ] && [ "$RESTART" = false ] && [ "$LOGS" = false ] && [ "$CLEAN" = false ]; then
        DEPLOY=true
    fi
    
    # 执行操作
    if [ "$CLEAN" = true ]; then
        clean_all
        exit 0
    fi
    
    if [ "$STOP" = true ]; then
        stop_existing
        exit 0
    fi
    
    if [ "$LOGS" = true ]; then
        show_logs
        exit 0
    fi
    
    if [ "$RESTART" = true ]; then
        restart_service
        exit 0
    fi
    
    if [ "$BUILD" = true ] || [ "$DEPLOY" = true ]; then
        check_dependencies
        create_directories
        
        if [ "$BUILD" = true ]; then
            build_image
        fi
        
        if [ "$DEPLOY" = true ]; then
            stop_existing
            
            if [ "$USE_COMPOSE" = true ]; then
                deploy_with_compose
            else
                deploy_with_docker
            fi
            
            check_deployment
            show_deployment_info
        fi
    fi
}

# 运行主函数
main "$@"
