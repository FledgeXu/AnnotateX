package api

import (
	"annotate-x/httperr"
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DatasetHandler struct {
	DatasetService service.IDatasetService
}

func RegisterDatasetRouters(rg *gin.RouterGroup, datasetService service.IDatasetService) {
	handler := &DatasetHandler{datasetService}
	group := rg.Group("/datasets")
	group.POST("/create", handler.create)
}

func (h *DatasetHandler) create(c *gin.Context) {
	var createDatasetForm models.CreateDatasetForm
	if err := c.ShouldBindWith(&createDatasetForm, binding.FormMultipart); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}

	exist, err := h.DatasetService.DetermineIsExist(c.Request.Context(), createDatasetForm.Name, createDatasetForm.ProjectId)
	if err != nil {
		c.Error(err)
		return
	}
	if exist {
		c.Error(httperr.NewBadRequestError("Dataset is exist."))
		return
	}

	if err = h.DatasetService.Create(c.Request.Context(), &createDatasetForm); err != nil {
		c.Error(err)
		return
	}
	utils.OK(c, "created")
}
