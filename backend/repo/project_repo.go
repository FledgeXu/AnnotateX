package repo

import (
	"annotate-x/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type IProjectRepo interface {
	CreateProject(ctx context.Context, projectReq *models.CreateProjectRequest) (*models.Project, error)
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

func (r *ProjectRepo) CreateProject(ctx context.Context, projectReq *models.CreateProjectRequest) (*models.Project, error) {
	var project *models.Project
	query := `
		INSERT INTO projects (name, modality, description)
		VALUES (:name, :modality, :description) 
		RETURNING id, name, modality, status, description, created_at, updated_at
	`
	err := r.DB.GetContext(ctx, &project, query, projectReq)
	return project, err
}
