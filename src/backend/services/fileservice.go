package services

import (
	"os"
	"path/filepath"
	"simple-server/src/backend/config"
	"simple-server/src/backend/utils"
	"strings"
)

type FileService struct {
	config *config.Config
}

type FileEntry struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
}

type SearchResult struct {
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	RelativePath string `json:"relativePath"`
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{
		config: cfg,
	}
}

// ListFiles lists files in a directory
func (fs *FileService) ListFiles(relativePath string) ([]FileEntry, error) {
	// Clean path
	safePath := utils.SanitizePath(relativePath)
	fullPath := filepath.Join(fs.config.Storage.UploadDir, safePath)

	// Validate path security
	if !utils.IsValidPath(fs.config.Storage.UploadDir, safePath) {
		return nil, os.ErrInvalid
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []FileEntry
	for _, entry := range entries {
		// Filter hidden files and blocked directories
		if utils.IsHiddenFile(entry.Name()) {
			continue
		}

		if utils.IsBlockedPath(entry.Name(), fs.config.Security.BlockedPaths) {
			continue
		}

		// Handle symlinks
		info, err := entry.Info()
		if err != nil {
			continue
		}

		isDir := info.IsDir()
		if info.Mode()&os.ModeSymlink != 0 {
			// If it's a symlink, check the target
			linkPath := filepath.Join(fullPath, entry.Name())
			target, err := os.Stat(linkPath)
			if err == nil {
				isDir = target.IsDir()
			}
		}

		files = append(files, FileEntry{
			Name:        entry.Name(),
			IsDirectory: isDir,
		})
	}

	return files, nil
}

// SearchFiles searches for files
func (fs *FileService) SearchFiles(query, directory string) ([]SearchResult, error) {
	safePath := utils.SanitizePath(directory)
	searchPath := filepath.Join(fs.config.Storage.UploadDir, safePath)

	// Validate path security
	if !utils.IsValidPath(fs.config.Storage.UploadDir, safePath) {
		return nil, os.ErrInvalid
	}

	var results []SearchResult

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore error, continue searching
		}

		// Skip hidden directories (e.g. .git), do not traverse their contents
		if info.IsDir() && utils.IsHiddenDirectory(info.Name()) {
			return filepath.SkipDir
		}

		// Skip the directory itself, but continue traversing contents
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if utils.IsHiddenFile(info.Name()) {
			return nil
		}

		// Check if in blocked paths
		relativePath, _ := filepath.Rel(fs.config.Storage.UploadDir, path)
		if utils.IsBlockedPath(relativePath, fs.config.Security.BlockedPaths) {
			return nil
		}

		// Check if filename matches query
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
			results = append(results, SearchResult{
				FileName:     info.Name(),
				FilePath:     path,
				RelativePath: relativePath,
			})
		}

		return nil
	})

	return results, err
}

// FileExists checks if a file exists
func (fs *FileService) FileExists(relativePath string) bool {
	safePath := utils.SanitizePath(relativePath)
	fullPath := filepath.Join(fs.config.Storage.UploadDir, safePath)

	if !utils.IsValidPath(fs.config.Storage.UploadDir, safePath) {
		return false
	}

	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFullPath gets the full path of a file
func (fs *FileService) GetFullPath(relativePath string) string {
	safePath := utils.SanitizePath(relativePath)
	return filepath.Join(fs.config.Storage.UploadDir, safePath)
}

// IsMarkdownFile checks if a file is a Markdown file
func (fs *FileService) IsMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".markdown"
}

// ReadMarkdownFile reads the content of a Markdown file
func (fs *FileService) ReadMarkdownFile(relativePath string) ([]byte, error) {
	safePath := utils.SanitizePath(relativePath)
	fullPath := filepath.Join(fs.config.Storage.UploadDir, safePath)

	if !utils.IsValidPath(fs.config.Storage.UploadDir, safePath) {
		return nil, os.ErrInvalid
	}

	// Use the original relative path to check file extension, not the full path
	if !fs.IsMarkdownFile(relativePath) {
		return nil, os.ErrInvalid
	}

	return os.ReadFile(fullPath)
}
