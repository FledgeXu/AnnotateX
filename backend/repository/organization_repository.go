package repository

import (
	"annotate-x/model"

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

func (r *OrganizationRepository) CreateOrganization(organization *model.Organization) error {
	err := r.DB.QueryRow(`
		INSERT INTO organization (type, name, code, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		organization.OrgType,
		organization.Name,
		organization.Code,
		organization.Description,
	).Scan(&organization.ID)
	return err
}
