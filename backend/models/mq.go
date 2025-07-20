package models

type TransformDatasetMessage struct {
	Dataset Dataset  `json:"dataset"`
	Keys    []string `json:"keys"`
}
