package model

import "time"

type Project struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Modality    string    `db:"modality"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
