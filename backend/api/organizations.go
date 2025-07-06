package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/utils"
	"net/http"
	"strconv"

	"annotate-x/internal/context"

	"github.com/gin-gonic/gin"
)

func RegisterOrganizationsRouters(rg *gin.RouterGroup) {
	group := rg.Group("/organizations")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))
	group.POST("/create", create)
	group.GET("/:id", info)
}

func create(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req model.OrganizationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	exists, err := appCtx.OrgRepo.OrganizationExists(req.Name)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	if exists {
		utils.Unauthorized(c, "Organization exists")
		return
	}

	organization := model.Organization{
		Type:        req.Type,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	if err := appCtx.OrgRepo.CreateOrganization(&organization); err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, organization)
}

func info(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid organization ID: must be a numeric value")
		return
	}
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	model, err := appCtx.OrgRepo.GetOrganizationById(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.OK(c, model)
}
