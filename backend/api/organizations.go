package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/utils"

	"annotate-x/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationsHandler struct {
	OrgRepo *repository.OrganizationRepository
}

func RegisterOrganizationsRouters(rg *gin.RouterGroup,
	orgRepo *repository.OrganizationRepository,
) {
	handler := &OrganizationsHandler{orgRepo}
	group := rg.Group("/organizations")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))
	group.POST("/create", handler.create)
	group.GET("/:id", handler.info)
}

func (h *OrganizationsHandler) create(c *gin.Context) {
	var req model.OrganizationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	exists, err := h.OrgRepo.OrganizationExists(req.Name)
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
	if err := h.OrgRepo.CreateOrganization(&organization); err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, organization)
}

func (h *OrganizationsHandler) info(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid organization ID: must be a numeric value")
		return
	}
	model, err := h.OrgRepo.GetOrganizationById(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.OK(c, model)
}
