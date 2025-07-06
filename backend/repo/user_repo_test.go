package repository_test

import (
	"annotate-x/model"
	"annotate-x/repo"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	user := &model.User{
		Username:    "testuser",
		Password:    "hashpass",
		DisplayName: "Test User",
		Email:       "test@example.com",
		IsActive:    true,
	}

	mock.ExpectPrepare("INSERT INTO users").
		ExpectQuery().
		WithArgs(user.Username, user.Password, user.DisplayName, user.Email, user.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestGetUserByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	mockUser := model.User{
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

	user, err := repo.GetUserByID(mockUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, user.Username)
}

func TestGetUserByUsername(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	mockUser := model.User{
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

	user, err := repo.GetUserByUsername(mockUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, user.ID)
}

func TestUsernameExists(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users WHERE username = \$1\)`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.UsernameExists("testuser")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUpdateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	user := &model.User{
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

	err := repo.UpdateUser(user)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewUserRepo(db)

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(1)
	assert.NoError(t, err)
}
