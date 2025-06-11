package model

type OrganizationCreateRequest struct {
	Type        string `json:"type" binding:"required,oneof=vendor internal partner"` // OrganizationType
	Name        string `json:"name" binding:"required"`                               // Organization name
	Code        string `json:"code" binding:"omitempty,alphanum"`                     // Optional short code or slug
	Description string `json:"description" binding:"omitempty"`                       // Optional description
}

type Organization struct {
	ID          int64  `db:"id"`
	OrgType     string `db:"type"`
	Name        string `db:"name"`
	Code        string `db:"code"`
	Description string `db:"description"`
	IsActive    bool   `db:"is_active"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}
