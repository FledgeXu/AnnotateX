package repository

import (
	"github.com/jmoiron/sqlx"
)

type OrganizationRepository struct {
	DB *sqlx.DB
}

func NewOrganizationRepository(db *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{DB: db}
}

func (r *OrganizationRepository) OrganizationExists(name string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE name = $1
		)
	`, name)
	return exists, err
}
