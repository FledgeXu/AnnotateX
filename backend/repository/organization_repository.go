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
			SELECT 1 FROM organizations WHERE name = $1
		)
	`, name)
	return exists, err
}

func (r *OrganizationRepository) GetOrganizationById(id int64) (model.Organization, error) {
	var organization model.Organization
	err := r.DB.Get(&organization, `SELECT * FROM organizations WHERE id = $1`, id)
	return organization, err
}

func (r *OrganizationRepository) CreateOrganization(organization *model.Organization) error {
	err := r.DB.QueryRow(`
		INSERT INTO organizations (type, name, code, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		organization.Type,
		organization.Name,
		organization.Code,
		organization.Description,
	).Scan(&organization.ID)
	return err
}
