package main

import (
	"annotate-x/config"
	"annotate-x/internal/auth"

	"annotate-x/router"
)

func initialize() {
	auth.InitCasbinEnforcer()
}

func main() {
	initialize()
	r := router.SetupRouter()
	r.Run(config.AppConfig.LISTEN_ADDRESS)
}
