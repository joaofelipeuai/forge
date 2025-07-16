package forge

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadConfig represents file upload configuration
type UploadConfig struct {
	MaxFileSize   int64    // Maximum file size in bytes
	AllowedTypes  []string // Allowed MIME types
	UploadDir     string   // Upload directory
	GenerateName  bool     // Generate unique filename
	PreserveExt   bool     // Preserve file extension when generating name
	CreateDirs    bool     // Create directories if they don't exist
}

// UploadedFile represents an uploaded file
type UploadedFile struct {
	OriginalName string
	Filename     string
	Size         int64
	ContentType  string
	Path         string
	URL          string
}

// UploadResult represents the result of file upload
type UploadResult struct {
	Files   []*UploadedFile
	Errors  []error
	Success bool
}

// NewUploadConfig creates a new upload configuration
func NewUploadConfig(uploadDir string) *UploadConfig {
	return &UploadConfig{
		MaxFileSize:  10 << 20, // 10MB default
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "text/plain", "application/pdf"},
		UploadDir:    uploadDir,
		GenerateName: true,
		PreserveExt:  true,
		CreateDirs:   true,
	}
}

// FileUpload middleware
func FileUpload(config *UploadConfig) MiddlewareFunc {
	return func(c *Context) error {
		// Create upload directory if it doesn't exist
		if config.CreateDirs {
			if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
				return fmt.Errorf("failed to create upload directory: %v", err)
			}
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(config.MaxFileSize); err != nil {
			return fmt.Errorf("failed to parse multipart form: %v", err)
		}

		// Process uploaded files
		result := &UploadResult{
			Files:  make([]*UploadedFile, 0),
			Errors: make([]error, 0),
		}

		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
			for fieldName, files := range c.Request.MultipartForm.File {
				for _, fileHeader := range files {
					uploadedFile, err := ProcessUploadedFile(fileHeader, config)
					if err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("field %s: %v", fieldName, err))
						continue
					}
					result.Files = append(result.Files, uploadedFile)
				}
			}
		}

		result.Success = len(result.Errors) == 0

		// Store result in context
		c.Set("upload_result", result)

		return c.Next()
	}
}

// ProcessUploadedFile processes a single uploaded file
func ProcessUploadedFile(fileHeader *multipart.FileHeader, config *UploadConfig) (*UploadedFile, error) {
	// Check file size
	if fileHeader.Size > config.MaxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", fileHeader.Size, config.MaxFileSize)
	}

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Check content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !IsAllowedType(contentType, config.AllowedTypes) {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Generate filename
	filename := fileHeader.Filename
	if config.GenerateName {
		filename = GenerateUniqueFilename(fileHeader.Filename, config.PreserveExt)
	}

	// Create destination file
	destPath := filepath.Join(config.UploadDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	size, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Create uploaded file info
	uploadedFile := &UploadedFile{
		OriginalName: fileHeader.Filename,
		Filename:     filename,
		Size:         size,
		ContentType:  contentType,
		Path:         destPath,
		URL:          "/uploads/" + filename, // Assuming /uploads/ is the URL prefix
	}

	return uploadedFile, nil
}

// IsAllowedType checks if the content type is allowed
func IsAllowedType(contentType string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return true // Allow all types if none specified
	}

	for _, allowed := range allowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
		// Support wildcard matching (e.g., "image/*")
		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "/*")
			if strings.HasPrefix(contentType, prefix+"/") {
				return true
			}
		}
	}
	return false
}

// GenerateUniqueFilename generates a unique filename
func GenerateUniqueFilename(originalName string, preserveExt bool) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	if preserveExt {
		ext := filepath.Ext(originalName)
		return fmt.Sprintf("%d_%s%s", timestamp, randomStr, ext)
	}

	return fmt.Sprintf("%d_%s", timestamp, randomStr)
}

// Context methods for file upload
func (c *Context) GetUploadResult() *UploadResult {
	if result := c.Get("upload_result"); result != nil {
		return result.(*UploadResult)
	}
	return nil
}

func (c *Context) GetUploadedFiles() []*UploadedFile {
	if result := c.GetUploadResult(); result != nil {
		return result.Files
	}
	return nil
}

func (c *Context) GetUploadErrors() []error {
	if result := c.GetUploadResult(); result != nil {
		return result.Errors
	}
	return nil
}

// Single file upload helper
func SingleFileUpload(config *UploadConfig, fieldName string) MiddlewareFunc {
	return func(c *Context) error {
		if err := c.Request.ParseMultipartForm(config.MaxFileSize); err != nil {
			return fmt.Errorf("failed to parse multipart form: %v", err)
		}

		file, fileHeader, err := c.Request.FormFile(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get file from field %s: %v", fieldName, err)
		}
		defer file.Close()

		// Create a multipart.FileHeader slice for processing
		uploadedFile, err := ProcessUploadedFile(fileHeader, config)
		if err != nil {
			return err
		}

		// Store single file in context
		c.Set("uploaded_file", uploadedFile)

		return c.Next()
	}
}

// Get single uploaded file
func (c *Context) GetUploadedFile() *UploadedFile {
	if file := c.Get("uploaded_file"); file != nil {
		return file.(*UploadedFile)
	}
	return nil
}

// Static file serving for uploaded files
func (f *Forge) ServeUploads(urlPrefix, uploadDir string) {
	f.GET(urlPrefix+"/*filepath", func(c *Context) error {
		filePath := c.Params["filepath"]
		if filePath == "" {
			return c.String(404, "File not found")
		}

		// Security: prevent directory traversal
		if strings.Contains(filePath, "..") {
			return c.String(403, "Forbidden")
		}

		fullPath := filepath.Join(uploadDir, filePath)
		
		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return c.String(404, "File not found")
		}

		// Serve file
		http.ServeFile(c.Response, c.Request, fullPath)
		return nil
	})
}

// Image upload with validation
func ImageUpload(uploadDir string, maxSize int64) MiddlewareFunc {
	config := NewUploadConfig(uploadDir)
	config.MaxFileSize = maxSize
	config.AllowedTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	
	return FileUpload(config)
}

// Document upload with validation
func DocumentUpload(uploadDir string, maxSize int64) MiddlewareFunc {
	config := NewUploadConfig(uploadDir)
	config.MaxFileSize = maxSize
	config.AllowedTypes = []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain",
	}
	
	return FileUpload(config)
}