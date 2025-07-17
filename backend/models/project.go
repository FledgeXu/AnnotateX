package models

import "time"

type ProjectType string

const (
	ProjectType2D    ProjectType = "2D"
	ProjectType3D    ProjectType = "3D"
	ProjectTypeAudio ProjectType = "audio"
	ProjectTypeText  ProjectType = "text"
)

var ValidProjectTypes = []ProjectType{
	ProjectType2D,
	ProjectType3D,
	ProjectTypeAudio,
	ProjectTypeText,
}

type Project struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Modality    string    `db:"modality" json:"modality"`
	Status      string    `db:"status" json:"status"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Modality    string `json:"modality" binding:"omitempty"`
	Description string `json:"description"`
}

type ProjectFilter struct {
	OrderBy string
	Order   string
	Limit   int
	Offset  int
}
