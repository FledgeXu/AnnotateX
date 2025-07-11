package repo

import (
	"annotate-x/models"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type IProjectRepo interface {
	CreateProject(ctx context.Context, projectReq *models.CreateProjectRequest) error
	GetProjectByID(ctx context.Context, id int64) (*models.Project, error)
	ListProjects(ctx context.Context, filter models.ProjectFilter) ([]*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	// DeleteProject(ctx context.Context, id int64) error
	ProjectNameExists(ctx context.Context, name string) (bool, error)
}

type ProjectRepo struct {
	DB *sqlx.DB
}

func NewProjectRepo(db *sqlx.DB) *ProjectRepo {
	return &ProjectRepo{DB: db}
}

func (r *ProjectRepo) CreateProject(ctx context.Context, projectReq *models.CreateProjectRequest) error {
	query := `
		INSERT INTO projects (name, modality, description)
		VALUES (:name, :modality, :description) 
		RETURNING id, name, modality, status, description, created_at, updated_at
	`
	result, err := r.DB.NamedExecContext(ctx, query, projectReq)
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
