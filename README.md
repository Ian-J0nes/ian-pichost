# Ian·PicHost - 轻量级图片托管服务

```
 ██▓ ▄▄▄       ███▄    █         ▄▄▄██▀▀▀ ▒█████   ███▄    █ ▓█████   ██████ 
▓██▒▒████▄     ██ ▀█   █           ▒██   ▒██▒  ██▒ ██ ▀█   █ ▓█   ▀ ▒██    ▒ 
▒██▒▒██  ▀█▄  ▓██  ▀█ ██▒          ░██   ▒██░  ██▒▓██  ▀█ ██▒▒███   ░ ▓██▄   
░██░░██▄▄▄▄██ ▓██▒  ▐▌██▒       ▓██▄██▓  ▒██   ██░▓██▒  ▐▌██▒▒▓█  ▄   ▒   ██▒
░██░ ▓█   ▓██▒▒██░   ▓██░        ▓███▒   ░ ████▓▒░▒██░   ▓██░░▒████▒▒██████▒▒
░▓   ▒▒   ▓▒█░░ ▒░   ▒ ▒         ▒▓▒▒░   ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░░ ▒░ ░▒ ▒▓▒ ▒ ░
 ▒ ░  ▒   ▒▒ ░░ ░░   ░ ▒░        ▒ ░▒░     ░ ▒ ▒░ ░ ░░   ░ ▒░ ░ ░  ░░ ░▒  ░ ░
 ▒ ░  ░   ▒      ░   ░ ░         ░ ░ ░   ░ ░ ░ ▒     ░   ░ ░    ░   ░  ░  ░  
 ░        ░  ░         ░         ░   ░       ░ ░           ░    ░  ░      ░  
```

一个基于 Go 语言开发的轻量级图片托管服务，具有简洁优雅的黑色主题界面，支持拖拽上传、密钥验证、图片管理等功能。

## ✨ 特性

- 🚀 **高性能**: 基于 Gin 框架，轻量快速
- 🔐 **安全验证**: 支持密钥验证，保护您的图床服务
- 📱 **响应式设计**: 适配桌面和移动设备
- 🎨 **优雅界面**: 简约黑色主题，视觉体验佳
- 📁 **智能存储**: 按日期自动分类存储，MD5去重
- 🖼️ **多格式支持**: 支持 JPG、PNG、GIF、WebP 格式
- 📋 **图片管理**: 查看上传历史，一键复制链接
- 🐳 **容器化**: 支持 Docker 部署

## 🛠️ 技术栈

- **后端**: Go 1.24+ + Gin
- **前端**: 原生 HTML/CSS/JavaScript
- **存储**: 本地文件系统
- **部署**: Docker / 二进制文件

## 📦 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/yourusername/image-host.git
cd image-host

# 构建并运行
docker build -t ian-pichost .
docker run -d -p 8080:8080 -v $(pwd)/uploads:/app/uploads ian-pichost
```

### 方式二：源码编译

```bash
# 克隆项目
git clone https://github.com/yourusername/image-host.git
cd image-host

# 安装依赖
go mod download

# 编译运行
go build -o ian-pichost main.go
./ian-pichost
```

### 方式三：直接运行

```bash
go run main.go
```

## ⚙️ 配置说明

在 `main.go` 文件顶部可以自定义以下配置：

```go
const (
    AuthKey = "your-key"          // 访问密钥
    ServerPort = ":8080"          // 服务器端口
    MaxFileSize = 10 << 20        // 最大文件大小 (10MB)
    UploadDir = "./uploads"       // 上传目录
    ReadTimeout = 30 * time.Second // 读取超时
    WriteTimeout = 30 * time.Second // 写入超时
    IdleTimeout = 60 * time.Second  // 空闲超时
)

var SupportedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
```

## 🚀 使用方法

1. **启动服务**: 运行程序后访问 `http://localhost:8080`
2. **密钥验证**: 输入配置的访问密钥
3. **上传图片**: 点击或拖拽图片到上传区域
4. **管理图片**: 查看上传历史，复制图片链接

## 📡 API 接口

### 上传图片
```bash
POST /upload
Headers: X-Auth-Key: your-key
Content-Type: multipart/form-data
Body: image file
```

### 获取图片列表
```bash
GET /list
Headers: X-Auth-Key: your-key
```

### 访问图片
```bash
GET /images/{date}/{filename}
```

## 🐳 Docker 部署

### 使用预构建镜像
```bash
docker run -d \
  --name ian-pichost \
  -p 8080:8080 \
  -v /path/to/uploads:/app/uploads \
  -e AUTH_KEY=your-secret-key \
  ian-pichost:latest
```

### 自定义构建
```bash
docker build -t ian-pichost .
docker run -d -p 8080:8080 -v $(pwd)/uploads:/app/uploads ian-pichost
```

## 📁 项目结构

```
image-host/
├── main.go              # 主程序文件
├── templates/           # HTML 模板
│   ├── index.html      # 主页面
│   └── error.html      # 错误页面
├── static/             # 静态资源
│   └── app.js         # 前端脚本
├── uploads/           # 图片存储目录
├── go.mod            # Go 模块文件
├── go.sum            # 依赖校验文件
├── Dockerfile        # Docker 构建文件
├── docker-compose.yml # Docker Compose 配置
└── README.md         # 项目说明
```

## 🔧 开发指南

### 环境要求
- Go 1.24+
- Git

### 本地开发
```bash
# 克隆项目
git clone https://github.com/yourusername/image-host.git
cd image-host

# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go
```

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 👨‍💻 作者

**Leon** - *Initial work*

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- 感谢所有贡献者的支持

## 📞 支持

如果您遇到任何问题或有建议，请：

1. 查看 [Issues](https://github.com/yourusername/image-host/issues)
2. 创建新的 Issue
3. 联系作者

---

<div align="center">
  <strong>🖼️ Ian·PicHost Service 📸</strong><br>
  <em>Powered by Go • Built with ❤️</em><br><br>
  <strong>Leon</strong>
</div>
