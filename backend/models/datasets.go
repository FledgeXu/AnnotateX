package models

import "mime/multipart"

type CreateDatasetForm struct {
	Name      string                  `form:"name" binding:"required"`
	ProjectId int64                   `form:"project_id" binding:"required"`
	Files     []*multipart.FileHeader `form:"files" binding:"required"`
}
