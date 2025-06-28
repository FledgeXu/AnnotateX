package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/context"
	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/utils"
	"slices"

	"github.com/gin-gonic/gin"
)

func RegisterProjectsRouters(rg *gin.RouterGroup) {
	group := rg.Group("/projects")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))

	group.GET("/list", listProject)
	group.POST("/create", createProject)
}

func listProject(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	projects, err := appCtx.ProjectRepo.ListProjects()
	if err != nil {
		utils.InternalServerError(c, err.Error())
	}
	utils.OK(c, gin.H{
		"results": projects,
	})
}

func createProject(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !slices.Contains(model.ValidProjectTypes, model.ProjectType(req.Modality)) {
		utils.BadRequest(c, "Invalid Modality")
		return
	}

	isProjectExist, err := appCtx.ProjectRepo.ProjectNameExists(req.Name)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if isProjectExist {
		utils.BadRequest(c, "Project already exists")
		return
	}

	project, err := appCtx.ProjectRepo.CreateProject(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, project)
}
