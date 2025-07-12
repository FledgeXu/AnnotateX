package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.ProjectService.CreateProject(c.Request.Context(), req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, gin.H{})
}

func (h *ProjectHandler) getProjectById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}
	project, err := h.ProjectService.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, project)
}

func (h *ProjectHandler) listProjects(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := models.ProjectFilter{
		Limit:  limit,
		Offset: offset,
	}

	users, err := h.ProjectService.ListProjects(c.Request.Context(), filter)
	total := len(users)
	if err != nil {
		utils.InternalServerError(c, "Failed to get users")
		return
	}

	utils.OK(c, gin.H{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"results": users,
	})
}
