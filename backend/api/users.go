package api

import (
	"annotate-x/httperr"
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.IUserService
}

func RegisterUserRouters(rg *gin.RouterGroup, userService service.IUserService) {
	handler := &UserHandler{userService}
	group := rg.Group("/users")
	group.POST("/create", handler.createUser)
	group.GET("/:id", handler.getUserById)
	group.GET("/me", handler.getMe)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var req *models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}
	if err := h.UserService.Create(c.Request.Context(), req); err != nil {
		c.Error(httperr.NewBadRequestError(err.Error()))
		return
	}
	utils.Created(c, gin.H{})
}

func (h *UserHandler) getUserById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(httperr.NewBadRequestError("Invalid user ID"))
		return
	}
	userResp, err := h.UserService.GetUserById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.OK(c, userResp)
}

func (h *UserHandler) getMe(c *gin.Context) {
	userId, err := strconv.ParseInt(c.GetHeader(models.XUserID), 10, 64)
	if err != nil {
		// TODO: Fix this
		utils.Unauthorized(c, "Invalid Bearer Token")
		return
	}

	userResp, err := h.UserService.GetUserById(c.Request.Context(), userId)
	if err != nil {
		c.Error(err)
		return
	}

	utils.OK(c, userResp)
}
