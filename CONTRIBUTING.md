# 贡献指南

感谢您对 Ian·PicHost 项目的关注！我们欢迎任何形式的贡献，包括但不限于：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复
- ✨ 添加新功能

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Git
- Docker (可选)

### 本地开发设置

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 Fork 本项目
   ```

2. **克隆项目**
   ```bash
   git clone https://github.com/yourusername/image-host.git
   cd image-host
   ```

3. **安装依赖**
   ```bash
   go mod download
   ```

4. **运行开发服务器**
   ```bash
   make dev
   # 或者
   go run main.go
   ```

5. **访问应用**
   ```
   http://localhost:8080
   ```

## 📋 开发流程

### 1. 创建分支

为您的更改创建一个新分支：

```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/your-bug-fix
```

分支命名规范：
- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 2. 进行更改

- 保持代码风格一致
- 添加必要的注释
- 确保代码通过测试
- 更新相关文档

### 3. 提交更改

使用清晰的提交信息：

```bash
git add .
git commit -m "feat: 添加新的图片压缩功能"
```

提交信息格式：
- `feat:` - 新功能
- `fix:` - Bug 修复
- `docs:` - 文档更新
- `style:` - 代码格式化
- `refactor:` - 代码重构
- `test:` - 测试相关
- `chore:` - 构建过程或辅助工具的变动

### 4. 推送分支

```bash
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

在 GitHub 上创建 Pull Request，并：

- 提供清晰的标题和描述
- 说明更改的原因和内容
- 链接相关的 Issue（如果有）
- 添加截图（如果是 UI 更改）

## 🧪 测试

### 运行测试

```bash
make test
# 或
go test -v ./...
```

### 代码覆盖率

```bash
make test-coverage
```

### 代码质量检查

```bash
# 格式化代码
make fmt

# 静态检查
make vet

# Lint 检查（需要安装 golangci-lint）
make lint
```

## 📝 代码规范

### Go 代码规范

1. **遵循 Go 官方代码规范**
   - 使用 `gofmt` 格式化代码
   - 遵循 Go 命名约定
   - 添加适当的注释

2. **函数和变量命名**
   ```go
   // 好的命名
   func uploadImage(file *multipart.FileHeader) error
   var maxFileSize int64
   
   // 避免的命名
   func upload(f *multipart.FileHeader) error
   var mfs int64
   ```

3. **错误处理**
   ```go
   // 总是检查错误
   if err != nil {
       log.Printf("上传失败: %v", err)
       return err
   }
   ```

### 前端代码规范

1. **HTML**
   - 使用语义化标签
   - 保持适当的缩进
   - 添加必要的注释

2. **CSS**
   - 使用有意义的类名
   - 保持样式的一致性
   - 优先使用 CSS 变量

3. **JavaScript**
   - 使用现代 ES6+ 语法
   - 添加适当的注释
   - 处理错误情况

## 🐛 报告 Bug

### Bug 报告模板

创建 Issue 时，请包含以下信息：

```markdown
**Bug 描述**
简要描述遇到的问题

**复现步骤**
1. 访问 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

**期望行为**
描述您期望发生的情况

**实际行为**
描述实际发生的情况

**环境信息**
- OS: [例如 Windows 10, macOS 12.0, Ubuntu 20.04]
- 浏览器: [例如 Chrome 96, Firefox 95]
- Go 版本: [例如 1.24.3]

**截图**
如果适用，添加截图来帮助解释问题

**附加信息**
添加任何其他相关信息
```

## 💡 功能建议

### 功能请求模板

```markdown
**功能描述**
简要描述您希望添加的功能

**问题背景**
描述这个功能要解决的问题

**解决方案**
描述您希望的解决方案

**替代方案**
描述您考虑过的其他解决方案

**附加信息**
添加任何其他相关信息或截图
```

## 📚 文档贡献

### 文档类型

- **README.md** - 项目介绍和快速开始
- **API 文档** - 接口说明
- **部署文档** - 部署指南
- **开发文档** - 开发指南

### 文档规范

- 使用清晰的标题结构
- 提供代码示例
- 添加必要的截图
- 保持内容的准确性

## 🎯 开发建议

### 性能优化

- 关注内存使用
- 优化文件 I/O 操作
- 减少不必要的网络请求

### 安全考虑

- 验证用户输入
- 防止路径遍历攻击
- 限制文件大小和类型

### 用户体验

- 提供清晰的错误信息
- 优化加载时间
- 确保响应式设计

## 🤝 社区准则

### 行为准则

- 尊重所有贡献者
- 保持友好和专业
- 欢迎新手参与
- 提供建设性的反馈

### 沟通方式

- 使用 GitHub Issues 讨论问题
- 在 Pull Request 中进行代码审查
- 保持讨论的相关性

## 📞 获取帮助

如果您在贡献过程中遇到任何问题，可以：

1. 查看现有的 Issues 和 Pull Requests
2. 创建新的 Issue 寻求帮助
3. 联系项目维护者

## 🙏 致谢

感谢所有为 Ian·PicHost 项目做出贡献的开发者！

---

再次感谢您的贡献！每一个贡献都让 Ian·PicHost 变得更好。
