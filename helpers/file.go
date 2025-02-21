package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

func DownloadFileFromURL(url, localFilePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", localFilePath)
	return nil
}
