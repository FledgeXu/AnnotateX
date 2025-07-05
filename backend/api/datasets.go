package api

import (
	"annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterDatasetRouters(rg *gin.RouterGroup) {
	group := rg.Group("/datasets")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(auth.Enforcer))

	group.POST("/upload", uploadDataset)
}

func uploadDataset(c *gin.Context) {
	batchID := c.PostForm("batch_id")
	fmt.Println(batchID)
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	fmt.Println(file)
}
