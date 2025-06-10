package repository

import (
	"annotate-x/model"

	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"database/sql"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, `
		SELECT u.*, r.name AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.id = $1
		LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, `
		SELECT u.*, r.name AS role
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.username = $1
		LIMIT 1
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

func (r *UserRepository) CreateUser(user *model.User) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert user
	err = tx.QueryRowx(`
		INSERT INTO users (username, password_hash, display_name, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		user.Username,
		user.Password,
		user.DisplayName,
		user.Email,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	// Process user role, default is "unassigned".
	roleName := user.Role
	if roleName == "" {
		roleName = "unassigned"
	}

	var roleID int64
	err = tx.Get(&roleID, `SELECT id FROM roles WHERE name = $1`, roleName)
	if err != nil {
		return err
	}

	// Insert binding of role.
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`, user.ID, roleID)
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

func (r *UserRepository) FindWithFilter(f model.UserFilter) ([]model.User, int, error) {
	var users []model.User
	args := []any{}
	where := []string{}
	argIdx := 1

	// Fuzzy match on username or email
	if f.Keyword != "" {
		where = append(where, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+f.Keyword+"%")
		argIdx++
	}

	// is_active string filter
	if f.IsActive == "true" || f.IsActive == "false" {
		where = append(where, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, f.IsActive == "true")
		argIdx++
	}

	// Compose WHERE clause
	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	// Validate sort fields
	allowedSortFields := map[string]bool{
		"created_at": true,
		"username":   true,
		"email":      true,
	}
	sortBy := "created_at"
	if allowedSortFields[f.SortBy] {
		sortBy = f.SortBy
	}

	order := "DESC"
	if f.Order == "asc" {
		order = "ASC"
	}

	// Add LIMIT and OFFSET
	args = append(args, f.Limit, f.Offset)

	query := fmt.Sprintf(`
		SELECT id, username, email, is_active, created_at
		FROM users
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, sortBy, order, argIdx, argIdx+1)

	err := r.DB.Select(&users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	// Total count without pagination
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users %s`, whereSQL)
	var total int
	err = r.DB.Get(&total, countQuery, args[:argIdx-1]...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	result, err := r.DB.Exec(`
		UPDATE users
		SET password_hash = $1,
			display_name = $2,
		    email = $3,
		    is_active = $4,
		    updated_at = NOW()
		WHERE id = $5
	`, user.Password, user.DisplayName, user.Email, user.IsActive, user.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	// retrun updated user
	return r.GetUserByID(user.ID)
}
