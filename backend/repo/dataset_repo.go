package repo

import (
	"annotate-x/models"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type IDatasetRepo interface {
	CreateDataset(ctx context.Context, datasetReq *models.CreateDatasetRequest) error
	ExistsByNameAndProjectID(ctx context.Context, name string, projectID int64) (bool, error)
}

type DatasetRepo struct {
	DB *sqlx.DB
}

func NewDatasetRepo(db *sqlx.DB) *DatasetRepo {
	return &DatasetRepo{DB: db}
}

func (r *DatasetRepo) CreateDataset(ctx context.Context, datasetReq *models.CreateDatasetRequest) error {
	query := `
		INSERT INTO projects (project_id, name, description, format_version)
		VALUES (:project_id, :name, :description, :format_version) 
		RETURNING id, project_id, name, description, format_version, status, created_at, updated_at
	`
	result, err := r.DB.NamedExecContext(ctx, query, datasetReq)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
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
