package service

import (
	"annotate-x/models"
	"annotate-x/repo"
	"context"
	"errors"
)

type IProjectService interface {
	CreateProject(ctx context.Context, req *models.CreateProjectRequest) error
	UpdateProject(ctx context.Context, project *models.Project) error
	// DeleteProject(ctx context.Context, id int64) error
	GetProjectByID(ctx context.Context, id int64) (*models.Project, error)
	ListProjects(ctx context.Context, filter models.ProjectFilter) ([]*models.Project, error)
}

type ProjectService struct {
	ProjectRepo repo.IProjectRepo
}

func NewProjectService(projectRepo repo.IProjectRepo) *ProjectService {
	return &ProjectService{
		ProjectRepo: projectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *models.CreateProjectRequest) error {
	exist, err := s.ProjectRepo.ProjectNameExists(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("Project with same name existed.")
	}
	return s.ProjectRepo.CreateProject(ctx, req)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id int64) (*models.Project, error) {
	return s.ProjectRepo.GetProjectByID(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, filter models.ProjectFilter) ([]*models.Project, error) {
	return s.ProjectRepo.ListProjects(ctx, filter)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *models.Project) error {
	return s.ProjectRepo.UpdateProject(ctx, project)
}
