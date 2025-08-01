package handlers

import (
	"net/http"
	"path/filepath"
	"simple-server/internal/services"
	"simple-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// ListFiles handles file list requests
func (h *FileHandler) ListFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := h.fileService.ListFiles(path)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to read directory", err.Error())
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{"files": files})
}

// GetMarkdownContent handles Markdown content requests
func (h *FileHandler) GetMarkdownContent(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		utils.SendError(c, http.StatusBadRequest, "Path parameter is required")
		return
	}

	// Check if it is a Markdown file
	if !h.fileService.IsMarkdownFile(filePath) {
		utils.SendError(c, http.StatusBadRequest, "File is not a markdown file")
		return
	}

	content, err := h.fileService.ReadMarkdownFile(filePath)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "File not found or not accessible", err.Error())
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{
		"content":  string(content),
		"filename": filepath.Base(filePath),
		"path":     filePath,
	})
}

// SearchFiles handles file search requests
func (h *FileHandler) SearchFiles(c *gin.Context) {
	query := c.Query("q")
	directory := c.Query("dir")

	if query == "" {
		utils.SendError(c, http.StatusBadRequest, "Search query is required")
		return
	}

	if directory == "" {
		directory = ""
	}

	results, err := h.fileService.SearchFiles(query, directory)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Search failed", err.Error())
		return
	}

	utils.SendJSON(c, http.StatusOK, gin.H{
		"query": gin.H{
			"keyword":   query,
			"directory": directory,
		},
		"results": results,
		"count":   len(results),
	})
}
