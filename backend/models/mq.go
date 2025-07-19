package models

type TransformDatasetMessage struct {
	Name      string   `json:"name"`
	ProjectId int64    `json:"project_id"`
	Type      string   `json:"type"`
	Urls      []string `json:"urls"`
}
