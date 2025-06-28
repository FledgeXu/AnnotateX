package model

import "time"

type OrganizationCreateRequest struct {
	Type        string `json:"type" binding:"required,oneof=vendor internal partner"` // OrganizationType
	Name        string `json:"name" binding:"required"`                               // Organization name
	Code        string `json:"code" binding:"omitempty,ascii"`                        // Optional short code or slug
	Description string `json:"description" binding:"omitempty"`                       // Optional description
}

type Organization struct {
	ID          int64     `db:"id" json:"id"`
	Type        string    `db:"type" json:"type"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Description string    `db:"description" json:"description"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
