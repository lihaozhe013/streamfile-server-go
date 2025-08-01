package utils

import (
	"path/filepath"
	"strings"
)

// SanitizePath cleans and validates the path to prevent directory traversal attacks
func SanitizePath(inputPath string) string {
	// Clean path
	cleaned := filepath.Clean(inputPath)

	// Remove leading ".." and "/"
	cleaned = strings.TrimPrefix(cleaned, "../")
	cleaned = strings.TrimPrefix(cleaned, "/")

	// Remove any path segments containing ".."
	parts := strings.Split(cleaned, string(filepath.Separator))
	var safeParts []string

	for _, part := range parts {
		if part != ".." && part != "." && part != "" {
			safeParts = append(safeParts, part)
		}
	}

	return strings.Join(safeParts, string(filepath.Separator))
}

// IsValidPath checks if the path is valid and secure
func IsValidPath(basePath, requestPath string) bool {
	cleanPath := SanitizePath(requestPath)
	fullPath := filepath.Join(basePath, cleanPath)

	// Ensure the path is within the base directory
	absBase, _ := filepath.Abs(basePath)
	absFull, _ := filepath.Abs(fullPath)

	return strings.HasPrefix(absFull, absBase)
}

// IsHiddenFile checks if the filename is a hidden file (starts with .)
func IsHiddenFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

// IsHiddenDirectory checks if the directory name is a hidden directory (starts with .)
func IsHiddenDirectory(dirname string) bool {
	return strings.HasPrefix(dirname, ".")
}

// IsBlockedPath checks if the path is in the blocked list
func IsBlockedPath(path string, blockedPaths []string) bool {
	for _, blocked := range blockedPaths {
		if strings.HasPrefix(path, blocked) {
			return true
		}
	}
	return false
}
