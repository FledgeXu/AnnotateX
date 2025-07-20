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
	DetermineIsExist(ctx context.Context, name string, projectId int64) (bool, error)
	Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error
}

type DatasetService struct {
	DatasetRepo repo.IDatasetRepo
	S3Repo      repo.IS3Repo
	MqRepo      repo.IMqRepo
}

func NewDatasetService(datasetRepo repo.IDatasetRepo, s3Repo repo.IS3Repo, mqRepo repo.IMqRepo) *DatasetService {
	return &DatasetService{datasetRepo, s3Repo, mqRepo}
}

func (s *DatasetService) DetermineIsExist(ctx context.Context, name string, projectId int64) (bool, error) {
	return s.DatasetRepo.ExistsByNameAndProjectID(ctx, name, projectId)
}

func (s *DatasetService) Create(ctx context.Context, createDatasetForm *models.CreateDatasetForm) error {
	// TODO: Make this configurable
	limit := 10

	// Create dataset
	dataset, err := s.DatasetRepo.CreateDataset(ctx, &models.CreateDatasetRequest{
		ProjectId:     createDatasetForm.ProjectId,
		Name:          createDatasetForm.Name,
		Description:   createDatasetForm.Description,
		FormatVersion: createDatasetForm.FormatVersion,
	})
	if err != nil {
		return err
	}

	// Save files that user uploaded.
	tempDir, filePaths, err := saveFiles(ctx, createDatasetForm.Files, limit)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Transfer those files to s3.
	savedFileNames, err := s.uploadFiles(ctx, createDatasetForm.ProjectId, createDatasetForm.Name, filePaths, limit)
	if err != nil {
		return err
	}

	// Enqueue process task
	err = s.enqueueDatasetProcessTask(dataset, savedFileNames)
	if err != nil {
		return err
	}

	return nil
}

func (s *DatasetService) enqueueDatasetProcessTask(dataset *models.Dataset, keys []string) error {
	if err := s.MqRepo.DeclareQueue("dataset.create", true); err != nil {
		return err
	}
	return s.MqRepo.Publish("", "dataset.create", models.TransformDatasetMessage{
		Dataset: *dataset,
		Keys:    keys,
	})
}

func (s *DatasetService) uploadFiles(ctx context.Context, projectId int64, projectName string, filePaths []string, limit int) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	savedFileNames := make([]string, len(filePaths))
	for i, filePath := range filePaths {
		i, filePath := i, filePath
		g.Go(func() error {
			objectName := filepath.Join(strconv.FormatInt(projectId, 10), projectName, filepath.Base(filePath))
			savedFileNames[i] = objectName
			return s.S3Repo.UploadFile(ctx, objectName, filePath)
		})
	}

	err := g.Wait()
	return savedFileNames, err
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
