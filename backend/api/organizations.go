package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/utils"
	"fmt"
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
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	exists, err := appCtx.OrgRepo.OrganizationExists(req.Name)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, err.Error())
		return
	}
	if exists {
		utils.JSONError(c, http.StatusUnauthorized, "Organization exists")
		return
	}

	organization := model.Organization{
		Type:        req.Type,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	if err := appCtx.OrgRepo.CreateOrganization(&organization); err != nil {
		utils.JSONError(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusCreated, organization)
}

func info(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid organization ID: must be a numeric value")
		return
	}
	fmt.Println(id)
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	model, err := appCtx.OrgRepo.GetOrganizationById(id)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, model)
}
