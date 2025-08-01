# Simple Server Go 重构开发方案

## 项目概述
将现有的 Node.js 简单文件服务器重构为 Go 版本，保持相同的功能和API接口，提升性能和简化部署。
目前重构已经完成

## 核心功能分析（基于现有 server.ts）

### 1. 文件服务功能
- **静态文件服务**: 提供文件下载和浏览
- **目录浏览**: 自动生成目录列表，支持 index.html
- **文件上传**: 多文件上传到指定目录
- **Markdown 预览**: 自动渲染 .md 文件
- **文件搜索**: 基于文件名的搜索功能（当前使用 Rust 原生模块）

### 2. 安全访问控制
- **公共文件** (`files/`): 完全公开访问
- **私有文件** (`files/private-files/`): 仅通过直接URL访问
- **临时文件** (`files/incoming/`): 完全隐藏，仅用于上传暂存

### 3. API 端点
- `GET /` - 主页
- `GET /files/*` - 文件浏览和下载
- `GET /private-files/*` - 私有文件直接访问
- `GET /api/list-files` - 目录文件列表API
- `GET /api/markdown-content` - Markdown内容API
- `GET /api/search_feat/file_name=:fileName/current_dir=*` - 文件搜索API
- `POST /upload` - 文件上传API

## Go 重构项目结构

```
go/
├── main.go                # 应用入口点
├── go.mod                 # Go模块定义
├── go.sum                 # 依赖锁定文件
├── internal/              # 内部包，不对外暴露
│   ├── config/           # 配置管理
│   │   └── config.go
│   ├── handlers/         # HTTP处理器
│   │   ├── files.go      # 文件相关处理
│   │   ├── upload.go     # 上传处理
│   │   ├── markdown.go   # Markdown处理
│   │   └── search.go     # 搜索功能
│   ├── middleware/       # 中间件
│   │   ├── security.go   # 安全检查
│   │   └── logging.go    # 日志记录
│   ├── services/         # 业务逻辑服务
│   │   ├── fileservice.go # 文件操作服务
│   │   └── searchservice.go # 搜索服务
│   └── utils/            # 工具函数
│       ├── path.go       # 路径处理
│       └── response.go   # 响应工具
├── public/               # 静态文件
├── files/                # 文件存储目录
│   ├── incoming/         # 上传暂存
│   ├── private-files/    # 私有文件
│   └── [用户文件]         # 公共文件
└── README.md            # Go版本说明文档
```

## 技术栈选择

### 核心依赖
```go
// HTTP 路由和服务器
github.com/gin-gonic/gin  // 高性能 HTTP 框架

// 文件操作和工具
github.com/gorilla/mux    // 备选路由器
mime                      // MIME类型检测

// 配置管理
github.com/spf13/viper    // 配置管理
github.com/spf13/cobra    // CLI工具（可选）

// 日志
github.com/sirupsen/logrus // 结构化日志

// 文件搜索（替代Rust模块）
path/filepath             // Go标准库文件路径操作
strings                   // 字符串匹配
```

## API 规格定义

### 1. 文件列表 API
```
GET /api/list-files?path={相对路径}
Response: {
  "files": [
    {
      "name": "文件名",
      "isDirectory": true/false
    }
  ]
}
```

### 2. Markdown 内容 API
```
GET /api/markdown-content?path={文件相对路径}
Response: {
  "content": "markdown内容",
  "filename": "文件名",
  "path": "相对路径"
}
```

### 3. 文件搜索 API
```
GET /api/search?q={搜索关键词}&dir={搜索目录}
Response: {
  "query": {
    "keyword": "搜索词",
    "directory": "搜索目录"
  },
  "results": [
    {
      "fileName": "文件名",
      "filePath": "完整路径",
      "relativePath": "相对路径"
    }
  ],
  "count": 结果数量
}
```

### 4. 文件上传 API
```
POST /upload
Content-Type: multipart/form-data
Response: {
  "message": "上传成功",
  "filename": "上传的文件名"
}
```

## 配置管理

### 环境变量
- `HOST`: 服务器监听地址 (默认: 0.0.0.0)
- `PORT`: 服务器端口 (默认: 80)
- `UPLOAD_DIR`: 上传目录 (默认: ./files)
- `LOG_LEVEL`: 日志级别 (默认: info)

## 安全考虑

### 路径安全
- 使用 `filepath.Clean()` 清理路径
- 验证路径不包含 `../` 等危险模式
- 限制访问范围在指定目录内

### 文件类型验证
- 基于扩展名和MIME类型双重验证
- 可配置的允许/禁止文件类型列表

### 访问控制
- incoming 目录完全隐藏
- private-files 仅直接URL访问
- 隐藏以 `.` 开头的文件
