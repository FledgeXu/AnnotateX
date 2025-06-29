package repository

import (
	"annotate-x/model"

	"github.com/jmoiron/sqlx"
)

type ProjectRepository struct {
	DB *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (r *ProjectRepository) GetProjectByID(id int64) (*model.Project, error) {
	var project model.Project
	err := r.DB.Get(&project, `SELECT * FROM projects WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) ProjectNameExists(name string) (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS (
		SELECT 1 FROM projects WHERE name = $1
	)`
	err := r.DB.Get(&exists, query, name)
	return exists, err
}

func (r *ProjectRepository) CreateProject(req *model.CreateProjectRequest) (*model.Project, error) {
	query := `
	INSERT INTO projects (name, modality, status, description)
	VALUES (:name, :modality, :status, :description)
	RETURNING *
	`

	args := map[string]any{
		"name":        req.Name,
		"modality":    req.Modality,
		"status":      "active",
		"description": req.Description,
	}

	var project model.Project
	stmt, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(&project, args); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) ListProjects() ([]model.Project, error) {
	var projects []model.Project
	err := r.DB.Select(&projects, `SELECT * FROM projects ORDER BY created_at DESC`)
	return projects, err
}
