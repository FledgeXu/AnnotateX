package api

import (
	"annotate-x/service"
	"annotate-x/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.IUserService
}

func RegisterUserRouters(rg *gin.RouterGroup, userService service.IUserService) {
	handler := &UserHandler{userService}
	group := rg.Group("/users")
	group.GET("/:id", handler.getUserById)
}

func (h *UserHandler) getUserById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}
	userResp, err := h.UserService.GetUserById(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	utils.OK(c, userResp)
}
