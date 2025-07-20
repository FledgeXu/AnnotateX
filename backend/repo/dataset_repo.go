package repo

import (
	"annotate-x/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type IDatasetRepo interface {
	CreateDataset(ctx context.Context, datasetReq *models.CreateDatasetRequest) (*models.Dataset, error)
	ExistsByNameAndProjectID(ctx context.Context, name string, projectID int64) (bool, error)
}

type DatasetRepo struct {
	DB *sqlx.DB
}

func NewDatasetRepo(db *sqlx.DB) *DatasetRepo {
	return &DatasetRepo{DB: db}
}

func (r *DatasetRepo) CreateDataset(ctx context.Context, datasetReq *models.CreateDatasetRequest) (*models.Dataset, error) {
	query := `
		INSERT INTO datasets (project_id, name, description, format_version)
		VALUES (:project_id, :name, :description, :format_version)
		RETURNING id, project_id, name, description, format_version, status, created_at, updated_at
	`
	var dataset models.Dataset
	stmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := stmt.GetContext(ctx, &dataset, datasetReq); err != nil {
		return nil, err
	}
	return &dataset, nil
}

func (r *DatasetRepo) ExistsByNameAndProjectID(ctx context.Context, name string, projectID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM datasets
			WHERE name = $1 AND project_id = $2
		)
	`
	err := r.DB.GetContext(ctx, &exists, query, name, projectID)
	return exists, err
}
