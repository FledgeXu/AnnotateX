package main

import (
	"annotate-x/config"
	"annotate-x/internal/auth"

	"annotate-x/router"
)

func main() {
	auth.InitCasbinEnforcer()
	r := router.SetupRouter()
	r.Run(config.AppConfig.LISTEN_ADDRESS)
}
