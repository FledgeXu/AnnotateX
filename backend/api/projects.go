package api

import (
	"annotate-x/httperr"
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListProjectsQuery struct {
	Limit   int    `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset  int    `form:"offset,default=0" binding:"gte=0"`
	OrderBy string `form:"order_by,default=created_at" binding:"oneof=created_at name"`
	Order   string `form:"order,default=desc" binding:"oneof=asc desc"`
}

type ProjectHandler struct {
	ProjectService service.IProjectService
}

func RegisterProjectRouters(rg *gin.RouterGroup, projectService service.IProjectService) {
	handler := &ProjectHandler{projectService}
	group := rg.Group("/projects")
	group.POST("/create", handler.createProject)
	group.GET("/:id", handler.getProjectById)
	group.GET("/list", handler.listProjects)
}

func (h *ProjectHandler) createProject(c *gin.Context) {
	var req *models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}
	if err := h.ProjectService.CreateProject(c.Request.Context(), req); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}
	utils.Created(c, gin.H{})
}

func (h *ProjectHandler) getProjectById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(httperr.NewBadRequestError("Invalid user ID"))
		return
	}
	project, err := h.ProjectService.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.OK(c, project)
}

func (h *ProjectHandler) listProjects(c *gin.Context) {
	var query ListProjectsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(httperr.NewBadRequestError("invalid query parameters: " + err.Error()))
		return
	}

	filter := models.ProjectFilter{
		OrderBy: query.OrderBy,
		Order:   query.Order,
		Limit:   query.Limit,
		Offset:  query.Offset,
	}

	projects, err := h.ProjectService.ListProjects(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	utils.OK(c, gin.H{
		"limit":   query.Limit,
		"offset":  query.Offset,
		"total":   len(projects),
		"results": projects,
	})
}
