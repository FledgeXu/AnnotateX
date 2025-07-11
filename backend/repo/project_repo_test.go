package repo_test

import (
	"annotate-x/models"
	"annotate-x/repo"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewProjectRepo(db)
	context := context.Background()

	projectReq := &models.CreateProjectRequest{
		Name:        "testproject",
		Modality:    "2D",
		Description: "test description",
	}

	mock.ExpectExec("INSERT INTO projects").
		WithArgs(projectReq.Name, projectReq.Modality, projectReq.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateProject(context, projectReq)
	assert.NoError(t, err)
}
