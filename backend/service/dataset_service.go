package service

import (
	"annotate-x/config"
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils"
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type IDatasetService interface {
	Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error
}

type DatasetService struct {
	S3Repo repo.IS3Repo
	MqRepo repo.IMqRepo
}

func NewDatasetService(s3Repo repo.IS3Repo, mqRepo repo.IMqRepo) *DatasetService {
	return &DatasetService{s3Repo, mqRepo}
}

func (s *DatasetService) Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error {
	// TODO: Make this configurable
	limit := 10

	tempDir, filePaths, err := saveFiles(ctx, createDatasetForm.Files, limit)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	for _, filePath := range filePaths {
		fp := filePath
		g.Go(func() error {
			objectName := filepath.Join(strconv.FormatInt(createDatasetForm.ProjectId, 10), createDatasetForm.Name, filepath.Base(filePath))
			return s.S3Repo.UploadFile(ctx, objectName, fp)

		})
	}

	return g.Wait()
}

func saveFiles(ctx context.Context, files []*multipart.FileHeader, limit int) (string, []string, error) {
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

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)

	savedFiles := make([]string, len(files))

	for i, file := range files {
		i, file := i, file

		g.Go(func() error {
			fileOutputPath := filepath.Join(dir, file.Filename)
			if err := utils.SaveUploadedFile(file, fileOutputPath); err != nil {
				return err
			}
			savedFiles[i] = fileOutputPath
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return "", nil, err
	}
	return dir, savedFiles, nil
}
