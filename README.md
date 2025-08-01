### Simple Server Go

A lightweight file server written in Go that provides file upload, download, Browse, and Markdown preview features.

-----

### Quick Start

#### Compile and Run

```bash
# run
go run .

# compile
go build

# run
./simple-server
```

#### Configuration

Create a `config.yaml` file (refer to `config.yaml.example`) or use environment variables:

```bash
# use environment variables
export HOST=0.0.0.0
export PORT=8000
./simple-server
```

-----

### Directory Structure

```
files/
├── incoming/          # Staging area for uploads (hidden from Browse)
├── private-files/     # Private files (accessible only via direct URL)
└── [public files...]  # Public files and directories
```

-----

### Development

```bash
# run in development mode
go run .

# generate cross-platform builds
GOOS=linux GOARCH=amd64 go build -o simple-server-linux
GOOS=windows GOARCH=amd64 go build -o simple-server-windows.exe
GOOS=darwin GOARCH=amd64 go build -o simple-server-macos
```