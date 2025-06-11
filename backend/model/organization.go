package model

type OrganizationCreateRequest struct {
	Type        string `json:"type" binding:"required,oneof=vendor internal partner"` // OrganizationType
	Name        string `json:"name" binding:"required"`                               // Organization name
	Code        string `json:"code" binding:"omitempty,alphanum"`                     // Optional short code or slug
	Description string `json:"description" binding:"omitempty"`                       // Optional description
}
