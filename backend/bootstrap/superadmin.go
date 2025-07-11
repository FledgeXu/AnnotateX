package bootstrap

import (
	"annotate-x/config"
	"annotate-x/db"
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils/security"
	"context"
)

func CreateSuperAdmin() {
	appConfig := config.AppConfig
	db := db.InitDB(config.GetConfig().DATABASE_URL)
	userRepo := repo.NewUserRepo(db)
	ctx := context.Background()

	if exists, err := userRepo.UsernameExists(ctx, config.GetConfig().SUPER_ADMIN_USERNAME); err != nil {
		panic(err.Error())
	} else if exists {
		return
	}
	hashedPassword, err := security.HashPassword(config.GetConfig().SUPER_ADMIN_PASSWORD)
	if err != nil {
		panic(err.Error())
	}
	user := &models.User{
		Username:    appConfig.SUPER_ADMIN_USERNAME,
		Password:    hashedPassword,
		DisplayName: "superadmin",
		Email:       "",
		IsActive:    true,
	}
	if err := userRepo.CreateUser(ctx, user); err != nil {
		panic(err.Error())
	}
}
