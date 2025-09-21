# StreamFile Server Go

Go版本的轻量级文件服务器，提供文件上传、下载、浏览、Markdown预览，以及增强的视频/音频播放器功能。

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

## 媒体播放器

系统内置基于 Video.js (CDN 引入) 的视频/音频播放器，自动匹配常见媒体文件扩展名：

视频：`.mp4`, `.webm`, `.ogv`, `.mov`, `.m4v`, `.mkv`, `.avi`

音频：`.mp3`, `.wav`, `.ogg`, `.m4a`, `.flac`, `.aac`

使用方式：
1. 将媒体文件放在 `files/` 目录（非隐藏/受限目录）下。
2. 浏览器访问 `/files/路径/文件名.mp4`。
3. 将自动进入播放器页面而不是直接下载。
4. 如需直接访问原始文件（浏览器默认行为/下载），可在 URL 后添加 `?raw=1`，例如：`/files/video/test.mp4?raw=1`。

桌面端快捷键：
* 左 / 右方向键：后退 / 前进 5 秒
* 空格：播放 / 暂停
* F：全屏切换

说明：
* 播放器通过内部将真实媒体流加上 `?raw=1` 访问，依旧支持断点续传（Range 请求）。
* 音频文件使用精简布局。
* 若需自定义主题或本地化 Video.js，可后续引入自定义 CSS/JS。

