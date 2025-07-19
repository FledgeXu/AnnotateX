package service

import (
	"annotate-x/config"
	"annotate-x/models"
	"annotate-x/utils"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type IDatasetService interface {
	Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error
}

type DatasetService struct {
}

func NewDatasetService() *DatasetService {
	return &DatasetService{}
}

func (s *DatasetService) Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error {
	tempDir, filePaths, err := saveFiles(createDatasetForm.Files)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	fmt.Println(filePaths)
	return nil
}

func saveFiles(files []*multipart.FileHeader) (string, []string, error) {
	tempDir := config.GetConfig().TEMP_DIR
	if tempDir != "" {
		if err := os.MkdirAll(tempDir, 0700); err != nil && !os.IsExist(err) {
			return "", nil, err
		}
	}

	dir, err := os.MkdirTemp(tempDir, uuid.New().String())
	if err != nil {
		return "", nil, err
	}

	savedFiles := []string{}
	for _, file := range files {
		fileOutputPath := filepath.Join(dir, file.Filename)
		if err := utils.SaveUploadedFile(file, fileOutputPath); err != nil {
			return "", nil, err
		}
		savedFiles = append(savedFiles, fileOutputPath)
	}

	return dir, savedFiles, nil
}
