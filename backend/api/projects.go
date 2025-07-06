package api

import (
	"annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/repository"
	"annotate-x/utils"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectsHandler struct {
	ProjectRepo *repository.ProjectRepository
}

func RegisterProjectsRouters(rg *gin.RouterGroup,
	projectRepo *repository.ProjectRepository,
	cacheRepo *repository.CacheRepository,
) {
	handler := &ProjectsHandler{projectRepo}
	group := rg.Group("/projects")
	group.Use(middleware.AuthMiddleware(cacheRepo), middleware.RequirePermissionMiddleware(auth.Enforcer))

	group.GET("/list", handler.list)
	group.POST("/create", handler.create)
	group.GET("/:id", handler.get)
}

func (h *ProjectsHandler) list(c *gin.Context) {
	projects, err := h.ProjectRepo.ListProjects()
	if err != nil {
		utils.InternalServerError(c, err.Error())
	}
	if projects == nil {
		projects = []model.Project{}
	}
	utils.OK(c, projects)
}

func (h *ProjectsHandler) create(c *gin.Context) {
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !slices.Contains(model.ValidProjectTypes, model.ProjectType(req.Modality)) {
		utils.BadRequest(c, "Invalid Modality")
		return
	}

	isProjectExist, err := h.ProjectRepo.ProjectNameExists(req.Name)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if isProjectExist {
		utils.BadRequest(c, "Project already exists")
		return
	}

	project, err := h.ProjectRepo.CreateProject(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, project)
}

func (h *ProjectsHandler) get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid project ID")
		return
	}

	project, err := h.ProjectRepo.GetProjectByID(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if project == nil {
		utils.NotFound(c, "Project not found")
		return
	}

	utils.OK(c, project)
}
