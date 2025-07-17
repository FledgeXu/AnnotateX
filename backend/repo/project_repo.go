package repo

import (
	"annotate-x/models"
	"context"
	"database/sql"
	"fmt"

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

func (r *ProjectRepo) GetProjectByID(ctx context.Context, id int64) (*models.Project, error) {
	var project models.Project
	err := r.DB.GetContext(ctx, &project, `
	SELECT * 
	From projects 
	WHERE id = $1
	LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) ListProjects(ctx context.Context, filter models.ProjectFilter) ([]*models.Project, error) {
	projects := []*models.Project{}

	query := fmt.Sprintf(`
		SELECT id, name, modality, status, description, created_at, updated_at
		FROM projects
		ORDER BY %s %s
		LIMIT %d OFFSET %d
	`, filter.OrderBy, filter.Order, filter.Limit, filter.Offset)
	fmt.Println(query)
	err := r.DB.SelectContext(ctx, &projects, query)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepo) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = :name,
			modality = :modality,
			status = :status,
			description = :description,
			updated_at = NOW()
		WHERE id = :id
	`
	result, err := r.DB.NamedExecContext(ctx, query, project)
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

func (r *ProjectRepo) ProjectNameExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.DB.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM projects WHERE name = $1)`, name)
	return exists, err
}
