package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"simple-server/internal/config"
	"simple-server/internal/handlers"
	"simple-server/internal/middleware"
	"simple-server/internal/services"
	"simple-server/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Set up logging
	logger := logrus.New()

	// Check if logging is enabled
	if !cfg.Logging.Enabled {
		// Disable log output
		logger.SetOutput(io.Discard)
	} else {
		// Set log level
		level, err := logrus.ParseLevel(cfg.Logging.Level)
		if err != nil {
			level = logrus.InfoLevel
		}
		logger.SetLevel(level)

		// Set log format
		if cfg.Logging.Format == "json" {
			logger.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp: true,
			})
		}

		// Set log output location
		if cfg.Logging.ToFile {
			setupLogFile(cfg, logger)
		}
	}

	// Ensure required directories exist
	createDirectories(cfg, logger)

	// Initialize services
	fileService := services.NewFileService(cfg)

	// Initialize handlers
	fileHandler := handlers.NewFileHandler(fileService)
	uploadHandler := handlers.NewUploadHandler(cfg, logger)

	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.SecurityMiddleware(cfg))
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Set up static file service
	setupStaticRoutes(router, cfg)

	// Set up API routes
	setupAPIRoutes(router, fileHandler, uploadHandler)

	// Set up file service routes
	setupFileRoutes(router, cfg)

	// Print startup info
	printStartupInfo(cfg, logger)

	// Start server
	server := &http.Server{
		Addr:         cfg.GetListenAddr(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	logger.Infof("Server starting on %s", cfg.GetListenAddr())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

// createDirectories creates required directories
func createDirectories(cfg *config.Config, logger *logrus.Logger) {
	dirs := []string{
		cfg.Storage.UploadDir,
		cfg.Storage.IncomingDir,
		cfg.Storage.PrivateDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Copy 404 page to incoming and private directories
	source404 := "./public/404-index.html"
	if _, err := os.Stat(source404); err == nil {
		copyFile(source404, filepath.Join(cfg.Storage.IncomingDir, "index.html"), logger)
		copyFile(source404, filepath.Join(cfg.Storage.PrivateDir, "index.html"), logger)
	}
}

// copyFile copies a file
func copyFile(src, dst string, logger *logrus.Logger) {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		srcFile, err := os.Open(src)
		if err != nil {
			logger.Warnf("Failed to open source file %s: %v", src, err)
			return
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			logger.Warnf("Failed to create destination file %s: %v", dst, err)
			return
		}
		defer dstFile.Close()

		// Use io.Copy instead of WriteTo for compatibility
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			logger.Warnf("Failed to copy file from %s to %s: %v", src, dst, err)
		}
	}
}

// setupStaticRoutes sets static file routes
func setupStaticRoutes(router *gin.Engine, cfg *config.Config) {
	// Static file service (public directory)
	router.Static("/public", "./public")
	router.StaticFile("/", "./public/index.html")

	// Direct access to private files
	router.Static("/private-files", cfg.Storage.PrivateDir)
}

// setupAPIRoutes sets API routes
func setupAPIRoutes(router *gin.Engine, fileHandler *handlers.FileHandler, uploadHandler *handlers.UploadHandler) {
	api := router.Group("/api")
	{
		api.GET("/list-files", fileHandler.ListFiles)
		api.GET("/markdown-content", fileHandler.GetMarkdownContent)
		api.GET("/search", fileHandler.SearchFiles)
	}

	// Upload route
	router.POST("/upload", uploadHandler.UploadFile)
}

// setupFileRoutes sets file access routes
func setupFileRoutes(router *gin.Engine, cfg *config.Config) {
	// File browsing and download
	router.GET("/files/*filepath", func(c *gin.Context) {
		filePath := c.Param("filepath")

		// Use path handling function from utils package
		cleanPath := utils.SanitizePath(filePath)
		fullPath := filepath.Join(cfg.Storage.UploadDir, cleanPath)

		// Get absolute path for security check
		absUploadDir, err := filepath.Abs(cfg.Storage.UploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}

		// Security check: ensure path is within upload directory
		if !strings.HasPrefix(absFullPath+string(filepath.Separator), absUploadDir+string(filepath.Separator)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Check if it is the incoming directory
		absIncomingDir, _ := filepath.Abs(cfg.Storage.IncomingDir)
		if strings.HasPrefix(absFullPath+string(filepath.Separator), absIncomingDir+string(filepath.Separator)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Check if file/directory exists
		info, err := os.Stat(fullPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		if info.IsDir() {
			// Check if index.html exists
			indexPath := filepath.Join(fullPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
				return
			}
			// Otherwise return file browser
			c.File("./public/file-browser.html")
			return
		}

		ext := strings.ToLower(filepath.Ext(fullPath))

		// Raw param forces direct file serving (used by media player)
		if c.Query("raw") == "1" {
			c.File(fullPath)
			return
		}

		// Markdown viewer
		if ext == ".md" {
			c.File("./public/markdown-viewer.html")
			return
		}

		// Media player (video / audio)
		if isMediaExtension(ext) {
			c.File("./public/video-player.html")
			return
		}

		// Directly serve other files
		c.File(fullPath)
	})

	// File browser homepage
	router.GET("/files", func(c *gin.Context) {
		c.File("./public/file-browser.html")
	})
}

// isMediaExtension checks if the extension is a supported audio/video type
func isMediaExtension(ext string) bool {
	switch ext {
	// Video
	case ".mp4", ".webm", ".ogv", ".mov", ".m4v", ".mkv", ".avi":
		return true
	// Audio
	case ".mp3", ".wav", ".ogg", ".m4a", ".flac", ".aac":
		return true
	default:
		return false
	}
}

// printStartupInfo prints startup information
func printStartupInfo(cfg *config.Config, logger *logrus.Logger) {
	// Get local IP
	localIP := getLocalIP()

	logger.Infof("=== Simple Server Go ===")
	logger.Infof("Upload Directory: %s", cfg.Storage.UploadDir)
	logger.Infof("Max Upload Size: %d MB", cfg.Storage.MaxUploadSize/(1024*1024))

	if cfg.Server.Host == "0.0.0.0" {
		logger.Infof("Server accessible at: http://%s:%d", localIP, cfg.Server.Port)
		logger.Infof("Local access: http://localhost:%d", cfg.Server.Port)
	} else {
		logger.Infof("Server accessible at: http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	}
}

// getLocalIP gets the local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// setupLogFile sets up log file output
func setupLogFile(cfg *config.Config, logger *logrus.Logger) {
	// Create log directory
	if err := os.MkdirAll(cfg.Logging.LogDir, 0755); err != nil {
		logger.Warnf("Failed to create log directory %s: %v", cfg.Logging.LogDir, err)
		return
	}

	// Generate log file name (by date)
	logFileName := filepath.Join(cfg.Logging.LogDir, "server.log")

	// Open log file
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Warnf("Failed to open log file %s: %v", logFileName, err)
		return
	}

	// Set log output to file
	logger.SetOutput(logFile)

	// Print one-time info to console
	logrus.Printf("Log output switched to file: %s", logFileName)
}
