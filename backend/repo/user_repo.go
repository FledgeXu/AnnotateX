package repo

import (
	"annotate-x/models"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type IUserRepo interface {
	// Create
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	// Read
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	// ListUsers() ([]model.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)

	// Update
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateUserPassword(ctx context.Context, id int64, newHash string) error

	// Delete
	DeleteUser(ctx context.Context, id int64) error
}

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (username, password_hash, display_name, email, is_active)
		VALUES (:username, :password_hash, :display_name, :email, :is_active) 
		RETURNING id
	`
	stmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	err = stmt.Get(&id, user)
	return id, err
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.DB.GetContext(ctx, &user, `
		SELECT id, username, password_hash, display_name, email, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.DB.GetContext(ctx, &user, `
		SELECT id, username, password_hash, display_name, email, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.DB.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username)
	return exists, err
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = :username,
		    password_hash = :password_hash,
		    display_name = :display_name,
		    email = :email,
		    is_active = :is_active,
		    updated_at = NOW()
		WHERE id = :id
	`
	result, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *UserRepo) UpdateUserPassword(ctx context.Context, userID int64, newHash string) error {
	result, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`, newHash, userID)
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

func (r *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
