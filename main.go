/*
*
* ██▓ ▄▄▄       ███▄    █         ▄▄▄██▀▀▀ ▒█████   ███▄    █ ▓█████   ██████
*▓██▒▒████▄     ██ ▀█   █           ▒██   ▒██▒  ██▒ ██ ▀█   █ ▓█   ▀ ▒██    ▒
*▒██▒▒██  ▀█▄  ▓██  ▀█ ██▒          ░██   ▒██░  ██▒▓██  ▀█ ██▒▒███   ░ ▓██▄
*░██░░██▄▄▄▄██ ▓██▒  ▐▌██▒       ▓██▄██▓  ▒██   ██░▓██▒  ▐▌██▒▒▓█  ▄   ▒   ██▒
*░██░ ▓█   ▓██▒▒██░   ▓██░        ▓███▒   ░ ████▓▒░▒██░   ▓██░░▒████▒▒██████▒▒
*░▓   ▒▒   ▓▒█░░ ▒░   ▒ ▒         ▒▓▒▒░   ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░░ ▒░ ░▒ ▒▓▒ ▒ ░
* ▒ ░  ▒   ▒▒ ░░ ░░   ░ ▒░        ▒ ░▒░     ░ ▒ ▒░ ░ ░░   ░ ▒░ ░ ░  ░░ ░▒  ░ ░
* ▒ ░  ░   ▒      ░   ░ ░         ░ ░ ░   ░ ░ ░ ▒     ░   ░ ░    ░   ░  ░  ░
* ░        ░  ░         ░         ░   ░       ░ ░           ░    ░  ░      ░
 */

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ================================
// 可自定义配置项
// ================================

const (
	AuthKey = "your-key"

	ServerPort = ":8080" // 服务器监听端口

	MaxFileSize = 10 << 20    // 最大文件大小 (10MB)
	UploadDir   = "./uploads" // 上传目录

	ReadTimeout  = 30 * time.Second // 读取超时时间 (0.5分钟)
	WriteTimeout = 30 * time.Second // 写入超时时间 (1分钟)
	IdleTimeout  = 60 * time.Second // 空闲超时时间 (1分钟)
)

var SupportedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

func getSupportedFormatsMessage() string {
	formats := make([]string, len(SupportedImageTypes))
	for i, ext := range SupportedImageTypes {
		formats[i] = strings.TrimPrefix(ext, ".")
	}
	return "只支持 " + strings.Join(formats, ", ") + " 格式"
}

// ================================
// 以下为程序核心代码
// ================================

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
	Size    int64  `json:"size,omitempty"`
}

type ImageInfo struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Date     string `json:"date"`
}

func main() {

	fmt.Println(`

 ██▓ ▄▄▄       ███▄    █         ▄▄▄██▀▀▀ ▒█████   ███▄    █ ▓█████   ██████ 
▓██▒▒████▄     ██ ▀█   █           ▒██   ▒██▒  ██▒ ██ ▀█   █ ▓█   ▀ ▒██    ▒ 
▒██▒▒██  ▀█▄  ▓██  ▀█ ██▒          ░██   ▒██░  ██▒▓██  ▀█ ██▒▒███   ░ ▓██▄   
░██░░██▄▄▄▄██ ▓██▒  ▐▌██▒       ▓██▄██▓  ▒██   ██░▓██▒  ▐▌██▒▒▓█  ▄   ▒   ██▒
░██░ ▓█   ▓██▒▒██░   ▓██░        ▓███▒   ░ ████▓▒░▒██░   ▓██░░▒████▒▒██████▒▒
░▓   ▒▒   ▓▒█░░ ▒░   ▒ ▒         ▒▓▒▒░   ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░░ ▒░ ░▒ ▒▓▒ ▒ ░
 ▒ ░  ▒   ▒▒ ░░ ░░   ░ ▒░        ▒ ░▒░     ░ ▒ ▒░ ░ ░░   ░ ▒░ ░ ░  ░░ ░▒  ░ ░
 ▒ ░  ░   ▒      ░   ░ ░         ░ ░ ░   ░ ░ ░ ▒     ░   ░ ░    ░   ░  ░  ░  
 ░        ░  ░         ░         ░   ░       ░ ░           ░    ░  ░      ░  

═══════════════════════════════════════════════════════════════════════════════
                           🖼️  Picture Host Service  📸
                          Powered by Go • Built with ❤️
═══════════════════════════════════════════════════════════════════════════════
	`)
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Fatal("创建上传目录失败:", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.MaxMultipartMemory = MaxFileSize

	r.Use(func(c *gin.Context) {
		if c.Request.ContentLength > MaxFileSize {
			c.JSON(413, Response{
				Success: false,
				Message: fmt.Sprintf("请求体太大，最大支持 %s", formatFileSize(MaxFileSize)),
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// 路由设置
	r.GET("/", indexHandler)
	r.POST("/upload", authMiddleware(), uploadHandler)
	r.GET("/list", authMiddleware(), listHandler)
	r.Static("/images", UploadDir)
	r.Static("/static", "./static")
	r.GET("/ping", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 使用自定义服务器配置
	srv := &http.Server{
		Addr:           ServerPort,
		Handler:        r,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		IdleTimeout:    IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Fatal(srv.ListenAndServe())
}

// 验证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("X-Auth-Key")
		if authKey != AuthKey {
			// 对于API请求返回JSON错误
			if c.Request.URL.Path == "/upload" || c.Request.URL.Path == "/list" {
				c.JSON(401, Response{
					Success: false,
					Message: "未授权访问，请提供正确的验证密钥",
				})
			} else {
				// 对于页面请求返回HTML错误页面
				c.HTML(401, "error.html", gin.H{
					"title":   "访问被拒绝",
					"message": "未授权访问，请提供正确的验证密钥",
				})
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

// 首页处理器
func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Nano·pic",
	})
}

// 上传处理器 - VPS环境优化版本
func uploadHandler(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	log.Printf("开始处理上传请求 from %s", clientIP)

	// 设置响应头，防止前端超时
	c.Header("Connection", "keep-alive")
	c.Header("Cache-Control", "no-cache")

	// 增加更详细的错误处理
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("获取文件失败 from %s: %v", clientIP, err)
		var message string
		if strings.Contains(err.Error(), "timeout") {
			message = "上传超时，请检查网络连接或尝试上传更小的文件"
		} else if strings.Contains(err.Error(), "too large") {
			message = fmt.Sprintf("文件太大，最大支持 %s", formatFileSize(MaxFileSize))
		} else {
			message = "获取文件失败，请重试"
		}
		c.JSON(400, Response{
			Success: false,
			Message: message,
		})
		return
	}
	defer file.Close()

	log.Printf("文件信息 from %s: %s, 大小: %s", clientIP, header.Filename, formatFileSize(header.Size))

	// 双重检查文件大小
	if header.Size > MaxFileSize {
		log.Printf("文件太大 from %s: %s", clientIP, formatFileSize(header.Size))
		c.JSON(400, Response{
			Success: false,
			Message: fmt.Sprintf("文件太大，最大支持 %s", formatFileSize(MaxFileSize)),
		})
		return
	}

	// 检查文件是否为空
	if header.Size == 0 {
		log.Printf("空文件 from %s", clientIP)
		c.JSON(400, Response{
			Success: false,
			Message: "文件为空，请选择有效的图片文件",
		})
		return
	}

	if !isValidImageType(header.Filename) {
		c.JSON(400, Response{
			Success: false,
			Message: getSupportedFormatsMessage(),
		})
		return
	}

	// 创建目录结构
	ext := getFileExtension(header.Filename)
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(UploadDir, dateDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		log.Printf("创建目录失败: %v", err)
		c.JSON(500, Response{
			Success: false,
			Message: "创建目录失败",
		})
		return
	}

	log.Printf("开始读取和计算哈希值 from %s...", clientIP)
	hashStart := time.Now()

	// 使用更高效的方式：边读边计算hash，边写入临时文件
	hasher := md5.New()
	tempFile, err := os.CreateTemp(fullDir, "upload_*.tmp")
	if err != nil {
		log.Printf("创建临时文件失败 from %s: %v", clientIP, err)
		c.JSON(500, Response{
			Success: false,
			Message: "服务器存储错误，请稍后重试",
		})
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name()) // 清理临时文件
	}()

	multiWriter := io.MultiWriter(hasher, tempFile)

	buffer := make([]byte, 128*1024)

	var totalRead int64
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			totalRead += int64(n)
			// 每读取1MB，打印一次日志
			if totalRead%(1024*1024) == 0 {
				log.Printf("已读取 %s from %s", formatFileSize(totalRead), clientIP)
			}

			_, writeErr := multiWriter.Write(buffer[:n])
			if writeErr != nil {
				log.Printf("写入文件失败 from %s: %v", clientIP, writeErr)
				c.JSON(500, Response{
					Success: false,
					Message: "文件写入失败，请重试",
				})
				return
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("读取文件失败 from %s: %v", clientIP, err)
			c.JSON(500, Response{
				Success: false,
				Message: "文件读取失败，请检查网络连接",
			})
			return
		}
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	log.Printf("哈希计算完成 from %s: %s, 耗时: %v", clientIP, hash, time.Since(hashStart))

	filename := hash + ext
	finalPath := filepath.Join(fullDir, filename)

	if _, err := os.Stat(finalPath); err == nil {
		log.Printf("文件已存在 from %s: %s", clientIP, filename)
		imageURL := fmt.Sprintf("http://%s/images/%s/%s",
			c.Request.Host, dateDir, filename)
		c.JSON(200, Response{
			Success: true,
			Message: "文件已存在，无需重复上传",
			URL:     imageURL,
			Size:    header.Size,
		})
		return
	}

	moveStart := time.Now()
	if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		log.Printf("重命名失败，使用复制 from %s: %v", clientIP, err)
		tempFile.Seek(0, 0) // 重置文件指针

		finalFile, err := os.Create(finalPath)
		if err != nil {
			log.Printf("创建最终文件失败 from %s: %v", clientIP, err)
			c.JSON(500, Response{
				Success: false,
				Message: "服务器存储错误，请稍后重试",
			})
			return
		}
		defer finalFile.Close()

		_, err = io.CopyBuffer(finalFile, tempFile, buffer)
		if err != nil {
			log.Printf("复制到最终文件失败 from %s: %v", clientIP, err)
			c.JSON(500, Response{
				Success: false,
				Message: "文件保存失败，请重试",
			})
			return
		}
	}
	log.Printf("文件保存完成 from %s，耗时: %v", clientIP, time.Since(moveStart))

	// 设置文件权限
	if err := os.Chmod(finalPath, 0644); err != nil {
		log.Printf("设置文件权限失败 from %s: %v", clientIP, err)
	}

	// 构建图片URL，支持HTTPS
	protocol := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/images/%s/%s",
		protocol, c.Request.Host, dateDir, filename)

	totalTime := time.Since(start)
	log.Printf("文件上传成功 from %s: %s (%s), 总耗时: %v",
		clientIP, filename, formatFileSize(header.Size), totalTime)

	c.JSON(200, Response{
		Success: true,
		Message: "上传成功",
		URL:     imageURL,
		Size:    header.Size,
	})
}

// 图片列表处理器
func listHandler(c *gin.Context) {
	images, err := getImageList()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, images)
}

// 获取图片列表
func getImageList() ([]ImageInfo, error) {
	var images []ImageInfo

	err := filepath.Walk(UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isValidImageType(info.Name()) {
			relPath, _ := filepath.Rel(UploadDir, path)
			relPath = strings.ReplaceAll(relPath, "\\", "/")

			imageInfo := ImageInfo{
				URL:      "/images/" + relPath,
				Filename: info.Name(),
				Size:     info.Size(),
				Date:     info.ModTime().Format("2006-01-02 15:04:05"),
			}
			images = append(images, imageInfo)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按时间倒序排序
	sort.Slice(images, func(i, j int) bool {
		return images[i].Date > images[j].Date
	})

	return images, nil
}

// 检查是否为有效的图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(getFileExtension(filename))

	for _, validExt := range SupportedImageTypes {
		if ext == validExt {
			return true
		}
	}
	return false
}

// 获取文件扩展名
func getFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
