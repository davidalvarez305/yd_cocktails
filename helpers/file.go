package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

func GetFileInfo(filePath string) (string, *multipart.FileHeader, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("error getting file info: %w", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("error reading file: %w", err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     info.Size(),
		// You can add more fields if necessary (like Content-Type)
	}

	return string(content), fileHeader, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
