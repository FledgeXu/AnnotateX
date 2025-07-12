package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	ProjectService service.IProjectService
}

func RegisterProjectRouters(rg *gin.RouterGroup, projectService service.IProjectService) {
	handler := &ProjectHandler{projectService}
	group := rg.Group("/projects")
	group.POST("/create", handler.createProject)
}

func (h *ProjectHandler) createProject(c *gin.Context) {
	var req *models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.ProjectService.CreateProject(c.Request.Context(), req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, gin.H{})
}
