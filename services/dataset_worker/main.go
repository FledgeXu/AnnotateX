package main

import (
	"dataset_worker/models"
	"fmt"
	"time"
)

func main() {
	message := models.TransformDatasetMessage{
		Dataset: models.Dataset{
			ID:            1,
			ProjectID:     1,
			Name:          "HH",
			Description:   nil,
			FormatVersion: "HHH",
			Status:        "HHH",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		Keys: []string{"/1/test2/final.zip"},
	}
	fmt.Println(message)
}
