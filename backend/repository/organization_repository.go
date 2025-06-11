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
