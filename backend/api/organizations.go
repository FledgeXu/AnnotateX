package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/model"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := appCtx.OrgRepo.OrganizationExists(req.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization exists"})
		return
	}

	organization := model.Organization{
		Type:        req.Type,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	appCtx.OrgRepo.CreateOrganization(&organization)
	c.JSON(http.StatusCreated, organization)
}

func info(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid organization id",
		})
		return
	}
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	model, err := appCtx.OrgRepo.GetOrganizationById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid organization id",
		})
		return
	}
	c.JSON(http.StatusOK, model)
}
