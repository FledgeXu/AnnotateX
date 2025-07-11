package repo_test

//
// import (
// 	"annotate-x/models"
// 	"annotate-x/repo"
// 	"context"
// 	"testing"
// 	"time"
//
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestCreateProject(t *testing.T) {
// 	db, mock := setupMockDB(t)
// 	repo := repo.NewProjectRepo(db)
// 	context := context.Background()
//
// 	projectReq := &models.CreateProjectRequest{
// 		Name:        "testproject",
// 		Modality:    "2D",
// 		Description: "test description",
// 	}
//
// 	mock.ExpectPrepare("INSERT INTO projects").
// 		ExpectQuery().
// 		WithArgs(projectReq.Name, projectReq.Modality, projectReq.Description).
// 		WillReturnRows(
// 			sqlmock.NewRows([]string{
// 				"id", "name", "modality", "status", "description", "created_at", "updated_at",
// 			}).AddRow(
// 				1, "Demo", "image", "active", "Test project", time.Now(), time.Now(),
// 			),
// 		)
//
// 	project, err := repo.CreateProject(context, projectReq)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(1), project.ID)
// }
