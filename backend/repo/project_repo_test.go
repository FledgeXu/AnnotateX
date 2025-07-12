package repo_test

import (
	"annotate-x/models"
	"annotate-x/repo"
	"context"
	"testing"
	"time"

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

func TestGetProjectByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewProjectRepo(db)

	expected := models.Project{
		ID:          1,
		Name:        "Test",
		Modality:    "image",
		Status:      "active",
		Description: "desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery("(?i)SELECT \\* FROM projects WHERE id = \\$1 LIMIT 1").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "modality", "status", "description", "created_at", "updated_at",
		}).AddRow(expected.ID, expected.Name, expected.Modality, expected.Status, expected.Description, expected.CreatedAt, expected.UpdatedAt))

	ctx := context.Background()
	project, err := repo.GetProjectByID(ctx, expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, project.ID)
}

func TestListProjects(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewProjectRepo(db)

	filter := models.ProjectFilter{Limit: 10, Offset: 0}

	mock.ExpectQuery(`SELECT id, name, modality, status, description, created_at, updated_at FROM projects ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(filter.Limit, filter.Offset).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "modality", "status", "description", "created_at", "updated_at",
		}).AddRow(int64(1), "Test", "image", "active", "desc", time.Now(), time.Now()))

	ctx := context.Background()
	projects, err := repo.ListProjects(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Test", projects[0].Name)
}

func TestUpdateProject(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewProjectRepo(db)

	project := &models.Project{
		ID:          1,
		Name:        "Updated",
		Modality:    "image",
		Status:      "inactive",
		Description: "updated desc",
	}

	mock.ExpectExec(`UPDATE projects SET name = .*?, modality = .*?, status = .*?, description = .*?, updated_at = NOW\(\) WHERE id = .*?`).
		WithArgs(
			project.Name,
			project.Modality,
			project.Status,
			project.Description,
			project.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := repo.UpdateProject(ctx, project)

	assert.NoError(t, err)
}

func TestProjectNameExists(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewProjectRepo(db)

	name := "TestProject"

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM projects WHERE name = \$1\)`).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ctx := context.Background()
	exists, err := repo.ProjectNameExists(ctx, name)

	assert.NoError(t, err)
	assert.True(t, exists)
}
