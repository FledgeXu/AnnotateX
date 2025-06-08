package main

import (
	"annotate-x/config"

	"annotate-x/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(config.AppConfig.LISTEN_ADDRESS)
}
