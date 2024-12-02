package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Allowed file extensions for images
var allowedExtensions = []string{".jpg", ".jpeg", ".png", ".gif"}

// UploadFile handles the file upload process
func UploadFile(file *multipart.FileHeader, uploadDir string) (string, error) {
	// Ensure the upload directory exists
	if err := createDirIfNotExist(uploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Open the uploaded file
	fileContent, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer fileContent.Close()

	// Validate file extension
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidExtension(fileExtension) {
		return "", errors.New("invalid file type, allowed types are .jpg, .jpeg, .png, .gif")
	}

	// Generate a unique file name
	fileName := generateUniqueFileName(file.Filename)

	// Create a new file on the server
	dstFile, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		return "", fmt.Errorf("failed to create file on server: %w", err)
	}
	defer dstFile.Close()

	// Copy the uploaded file content to the destination file
	_, err = io.Copy(dstFile, fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	// Return the file path
	return filepath.Join(uploadDir, fileName), nil
}

// createDirIfNotExist checks if a directory exists, and creates it if not
func createDirIfNotExist(dir string) error {
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// isValidExtension checks if the file extension is valid
func isValidExtension(extension string) bool {
	for _, ext := range allowedExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

// generateUniqueFileName generates a unique file name using a timestamp and the original file name
func generateUniqueFileName(originalFileName string) string {
	// Get the file extension
	ext := filepath.Ext(originalFileName)

	// Generate a unique file name using the timestamp and the file extension
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}
