package repository

import (
	"github.com/jmoiron/sqlx"

	"database/sql"
)

type User struct {
	ID          int64  `db:"id"`
	Username    string `db:"username"`
	Password    string `db:"password_hash"`
	DisplayName string `db:"display_name"`
	Email       string `db:"email"`
	AvatarURL   string `db:"avatar_url"`
	IsActive    bool   `db:"is_active"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByID(id int64) (*User, error) {
	var user User
	err := r.DB.Get(&user, `
		SELECT * FROM users WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.DB.Get(&user, `
		SELECT * FROM users WHERE username = $1
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1
		)
	`, username)
	return exists, err
}

func (r *UserRepository) CreateUser(user *User) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowx(`
		INSERT INTO users (username, password_hash, display_name, email, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		user.Username,
		user.Password,
		user.DisplayName,
		user.Email,
		user.AvatarURL,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	var unassignedRoleID int
	err = tx.Get(&unassignedRoleID, `SELECT id FROM roles WHERE name = 'unassigned'`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)
	`, user.ID, unassignedRoleID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) AssignUserRole(userID, roleID int64) error {
	_, err := r.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`, userID, roleID)
	return err
}

func (r *UserRepository) GetUserRoles(userID int64) ([]string, error) {
	var roles []string
	err := r.DB.Select(&roles, `
		SELECT r.name FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	return roles, err
}

func (r *UserRepository) UpdateUserPassword(userID int64, newHash string) error {
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
