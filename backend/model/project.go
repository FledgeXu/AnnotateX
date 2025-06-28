package model

import "time"

type ProjectType string

const (
	ProjectType2D    ProjectType = "2D"
	ProjectType3D    ProjectType = "3D"
	ProjectTypeAudio ProjectType = "audio"
	ProjectTypeText  ProjectType = "text"
)

var validProjectTypes = []ProjectType{
	ProjectType2D,
	ProjectType3D,
	ProjectTypeAudio,
	ProjectTypeText,
}

type Project struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Modality    string    `db:"modality"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
