package main

import (
	"annotate-x/config"
	"annotate-x/router"

	"github.com/gin-contrib/graceful"
)

func main() {
	router := router.SetupRouter()

	g, err := graceful.New(router)
	if err != nil {
		panic(err)
	}

	if err := g.Run(config.AppConfig.LISTEN_ADDRESS); err != nil {
		panic(err)
	}
}
