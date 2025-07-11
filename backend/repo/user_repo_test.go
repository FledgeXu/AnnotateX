package repo_test

import (
	"annotate-x/models"
	"annotate-x/repo"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	user := &models.User{
		Username:    "testuser",
		Password:    "hashpass",
		DisplayName: "Test User",
		Email:       "test@example.com",
		IsActive:    true,
	}

	mock.ExpectExec("^INSERT INTO users").
		WithArgs(user.Username, user.Password, user.DisplayName, user.Email, user.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateUser(context, user)
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	mockUser := models.User{
		ID:          1,
		Username:    "testuser",
		Password:    "hashpass",
		DisplayName: "Test User",
		Email:       "test@example.com",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "display_name", "email", "is_active", "created_at", "updated_at",
	}).AddRow(
		mockUser.ID, mockUser.Username, mockUser.Password, mockUser.DisplayName,
		mockUser.Email, mockUser.IsActive, mockUser.CreatedAt, mockUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, username, password_hash.*FROM users").
		WithArgs(mockUser.ID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(context, mockUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, user.Username)
}

func TestGetUserByUsername(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	mockUser := models.User{
		ID:          2,
		Username:    "byusername",
		Password:    "hashpass2",
		DisplayName: "User 2",
		Email:       "user2@example.com",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "display_name", "email", "is_active", "created_at", "updated_at",
	}).AddRow(
		mockUser.ID, mockUser.Username, mockUser.Password, mockUser.DisplayName,
		mockUser.Email, mockUser.IsActive, mockUser.CreatedAt, mockUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, username, password_hash.*FROM users").
		WithArgs(mockUser.Username).
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername(context, mockUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, user.ID)
}

func TestUsernameExists(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users WHERE username = \$1\)`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.UsernameExists(context, "testuser")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUpdateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	user := &models.User{
		ID:          1,
		Username:    "updated",
		Password:    "newhash",
		DisplayName: "Updated Name",
		Email:       "updated@example.com",
		IsActive:    true,
	}

	mock.ExpectExec(`UPDATE users.*`).
		WithArgs(
			user.Username,
			user.Password,
			user.DisplayName,
			user.Email,
			user.IsActive,
			user.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUser(context, user)
	assert.NoError(t, err)
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	userRepo := repo.NewUserRepo(db)
	context := context.Background()

	const (
		userID  = int64(123)
		newHash = "new_hash"
	)

	mock.ExpectExec(`UPDATE users`).
		WithArgs(newHash, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.UpdateUserPassword(context, userID, newHash)
	assert.NoError(t, err)

	mock.ExpectExec(`UPDATE users`).
		WithArgs(newHash, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = userRepo.UpdateUserPassword(context, userID, newHash)
	assert.Equal(t, sql.ErrNoRows, err)

	mock.ExpectExec(`UPDATE users`).
		WithArgs(newHash, userID).
		WillReturnError(sql.ErrConnDone)

	err = userRepo.UpdateUserPassword(context, userID, newHash)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repo.NewUserRepo(db)
	context := context.Background()

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(context, 1)
	assert.NoError(t, err)
}
