package api

import (
	"annotate-x/config"
	"annotate-x/httperr"
	"annotate-x/utils"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type CreateDatasetForm struct {
	Name      string                  `form:"name" binding:"required"`
	ProjectId int64                   `form:"project_id" binding:"required"`
	Files     []*multipart.FileHeader `form:"files" binding:"required"`
}

type DatasetHandler struct {
}

func RegisterDatasetRouters(rg *gin.RouterGroup) {
	handler := &DatasetHandler{}
	group := rg.Group("/datasets")
	group.POST("/create", handler.create)
}

func (h *DatasetHandler) create(c *gin.Context) {
	var createDatasetForm CreateDatasetForm
	if err := c.ShouldBindWith(&createDatasetForm, binding.FormMultipart); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}

	tempDir := config.GetConfig().TEMP_DIR
	if tempDir != "" {
		if err := os.MkdirAll(tempDir, 0700); err != nil && !os.IsExist(err) {
			c.Error(err)
			return
		}
	}

	dir, err := os.MkdirTemp(tempDir, uuid.New().String())
	if err != nil {
		c.Error(err)
		return
	}
	defer os.RemoveAll(dir)

	for _, file := range createDatasetForm.Files {
		c.SaveUploadedFile(file, filepath.Join(dir, file.Filename))
	}

	utils.OK(c, "created")
}
