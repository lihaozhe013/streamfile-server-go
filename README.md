# StreamFile Server Go

A lightweight streamfile server written in Go that provides file upload, download, browse, Markdown preview, and an enhanced media (video/audio) player.

-----

# Quick Start

## Compile and Run
### Build
```
# this will build both backend and frontend
npm run build
```

### Run
```
./simple-server
```

### Backend Only
```bash
# run
go run .

# compile
go build

# run
./simple-server
```

## Configuration

Create a `config.yaml` file (refer to `config.yaml.example`) or use environment variables:

```bash
# use environment variables
export HOST=0.0.0.0
export PORT=8000
./simple-server
```

-----

# Directory Structure

```
files/
├── incoming/          # Staging area for uploads (hidden from Browse)
├── private-files/     # Private files (accessible only via direct URL)
└── [public files...]  # Public files and directories
```

-----

# Development

```bash
# run in development mode
go run .

# generate cross-platform builds
GOOS=linux GOARCH=amd64 go build -o simple-server-linux
GOOS=windows GOARCH=amd64 go build -o simple-server-windows.exe
GOOS=darwin GOARCH=amd64 go build -o simple-server-macos
```

-----

# Media Player

The server includes an integrated Video/Audio player (powered by Video.js via CDN) for common media files:

Supported extensions (auto-detected):

Video: `.mp4`, `.webm`, `.ogv`, `.mov`, `.m4v`, `.mkv`, `.avi`

Audio: `.mp3`, `.wav`, `.ogg`, `.m4a`, `.flac`, `.aac`

Usage:
1. Place media files anywhere under `files/` (except hidden / restricted directories).
2. Visit `/files/path/to/movie.mp4` in the browser.
3. The media player page will load automatically instead of directly downloading the file.
4. To force the raw file (default browser handling or download), append `?raw=1` to the URL, e.g. `/files/path/to/movie.mp4?raw=1`.

Keyboard shortcuts (desktop):
* Left / Right Arrow: seek -/+5 seconds
* Space: play / pause
* F: toggle fullscreen

Notes:
* The player uses a CDN copy of Video.js (no extra install needed).
* Audio files get a simplified responsive layout.
* Range requests still work because the actual media stream is fetched with `?raw=1` internally.
