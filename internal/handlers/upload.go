package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"simple-server/internal/config"
	"simple-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UploadHandler struct {
	config *config.Config
	logger *logrus.Logger
}

func NewUploadHandler(cfg *config.Config, logger *logrus.Logger) *UploadHandler {
	return &UploadHandler{
		config: cfg,
		logger: logger,
	}
}

// UploadFile handles file upload
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "No file uploaded", err.Error())
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > h.config.Storage.MaxUploadSize {
		utils.SendError(c, http.StatusRequestEntityTooLarge, "File too large")
		return
	}

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if !h.isAllowedExtension(ext) {
		utils.SendError(c, http.StatusBadRequest, "File type not allowed")
		return
	}

	// Ensure incoming directory exists
	if err := os.MkdirAll(h.config.Storage.IncomingDir, 0755); err != nil {
		h.logger.WithError(err).Error("Failed to create incoming directory")
		utils.SendError(c, http.StatusInternalServerError, "Server error")
		return
	}

	// Clean filename, support UTF-8
	filename := header.Filename

	// Create destination file
	destPath := filepath.Join(h.config.Storage.IncomingDir, filename)
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create destination file")
		utils.SendError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer destFile.Close()

	// Copy file content
	_, err = io.Copy(destFile, file)
	if err != nil {
		h.logger.WithError(err).Error("Failed to copy file content")
		utils.SendError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	h.logger.WithFields(logrus.Fields{
		"filename": filename,
		"size":     header.Size,
	}).Info("File uploaded successfully")

	utils.SendSuccess(c, "File uploaded successfully", gin.H{
		"filename": filename,
		"size":     header.Size,
	})
}

// isAllowedExtension checks if the file extension is allowed
func (h *UploadHandler) isAllowedExtension(ext string) bool {
	if len(h.config.Security.AllowedExtensions) == 0 {
		return true // If no restriction, allow all extensions
	}

	for _, allowed := range h.config.Security.AllowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}
