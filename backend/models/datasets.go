package models

import (
	"mime/multipart"
	"time"
)

type Dataset struct {
	ID            int64     `json:"id" db:"id"`
	ProjectID     int64     `json:"project_id" db:"project_id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description" db:"description"`
	FormatVersion string    `json:"format_version" db:"format_version"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDatasetForm struct {
	ProjectId     int64                   `form:"project_id" binding:"required"`
	Name          string                  `form:"name" binding:"required"`
	Description   string                  `form:"description" binding:"required"`
	FormatVersion string                  `form:"description" binding:"required"`
	Files         []*multipart.FileHeader `form:"files" binding:"required"`
}

type CreateDatasetRequest struct {
	ProjectId     int64  `db:"project_id"`
	Name          string `db:"name"`
	Description   string `db:"description"`
	FormatVersion string `db:"format_version"`
}
