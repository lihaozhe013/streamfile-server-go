# Simple Server Go

Go版本的轻量级文件服务器，提供文件上传、下载、浏览和Markdown预览功能。

## 快速开始

### 编译和运行

```bash
# run
go run .

# 编译
go build

# 运行
./simple-server
```

### 配置

创建 `config.yaml` 文件（参考 `config.yaml.example`）或使用环境变量：

```bash
# 使用环境变量
export HOST=0.0.0.0
export PORT=8000
./simple-server
```

## 目录结构

```
files/
├── incoming/          # 上传暂存区（隐藏访问）
├── private-files/     # 私有文件（仅直接URL访问）
└── [公共文件...]      # 公共文件和目录
```

## 开发

```bash
# 开发模式运行
go run .

# 生成跨平台编译
GOOS=linux GOARCH=amd64 go build -o simple-server-linux
GOOS=windows GOARCH=amd64 go build -o simple-server-windows.exe
GOOS=darwin GOARCH=amd64 go build -o simple-server-macos
```
