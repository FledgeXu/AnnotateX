package utils

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func DetectContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	detected := http.DetectContentType(buf[:n])
	if detected != "application/octet-stream" {
		return detected, nil
	}

	ext := filepath.Ext(filePath)
	return mime.TypeByExtension(ext), nil
}
