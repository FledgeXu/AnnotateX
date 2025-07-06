package repo

import (
	"annotate-x/model"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	// Create
	CreateUser(user *model.User) (int64, error)
	// Read
	GetUserByID(id int64) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	// ListUsers() ([]model.User, error)
	UsernameExists(username string) (bool, error)

	// Update
	UpdateUser(user *model.User) error
	UpdateUserPassword(id int64, newHash string)

	// Delete
	DeleteUser(id int64) error
}

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(user *model.User) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (username, password_hash, display_name, email, is_active)
		VALUES (:username, :password_hash, :display_name, :email, :is_active) 
		RETURNING id
	`
	stmt, err := r.DB.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	err = stmt.Get(&id, user)
	return id, err
}

func (r *UserRepo) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, `
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

func (r *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, `
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

func (r *UserRepo) UsernameExists(username string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username)
	return exists, err
}

func (r *UserRepo) UpdateUser(user *model.User) error {
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
	result, err := r.DB.NamedExec(query, user)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *UserRepo) UpdateUserPassword(userID int64, newHash string) error {
	result, err := r.DB.Exec(`
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

func (r *UserRepo) DeleteUser(id int64) error {
	result, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
