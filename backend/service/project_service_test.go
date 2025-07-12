package service_test

import (
	"annotate-x/mocks"
	"annotate-x/models"
	"annotate-x/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_Create_Success(t *testing.T) {
	projectRepo := mocks.NewMockIProjectRepo(t)
	context := context.Background()

	projectRepo.On("ProjectNameExists", mock.Anything, "test").Return(false, nil)
	projectRepo.On("CreateProject", mock.Anything, mock.AnythingOfType("*models.CreateProjectRequest")).Return(nil)

	projectService := service.NewProjectService(projectRepo)

	req := &models.CreateProjectRequest{
		Name:        "test",
		Modality:    "2D",
		Description: "test Description",
	}

	err := projectService.CreateProject(context, req)

	assert.NoError(t, err)
	projectRepo.AssertExpectations(t)
}

func TestProjectService_CreateWithExistedName_FAILURE(t *testing.T) {
	projectRepo := mocks.NewMockIProjectRepo(t)
	context := context.Background()

	projectRepo.On("ProjectNameExists", mock.Anything, "test").Return(true, nil)

	projectService := service.NewProjectService(projectRepo)

	req := &models.CreateProjectRequest{
		Name:        "test",
		Modality:    "2D",
		Description: "test Description",
	}

	err := projectService.CreateProject(context, req)

	assert.Error(t, err)
	assert.Equal(t, "Project with same name existed.", err.Error())
	projectRepo.AssertExpectations(t)
}
