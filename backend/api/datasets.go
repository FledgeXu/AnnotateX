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

	if error := h.DatasetService.Create(c.Request.Context(), &createDatasetForm); error != nil {
		c.Error(error)
		return
	}
	utils.OK(c, "created")
}
