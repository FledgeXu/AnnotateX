package models

import "time"

type Dataset struct {
	ID            int64     `json:"id"`
	ProjectID     int64     `json:"project_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	FormatVersion string    `json:"format_version"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransformDatasetMessage struct {
	Dataset Dataset  `json:"dataset"`
	Keys    []string `json:"keys"`
}
