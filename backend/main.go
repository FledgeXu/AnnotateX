package main

import (
	"annotate-x/config"
	"annotate-x/db"
	"annotate-x/models"
	"annotate-x/mq"
	"annotate-x/router"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/graceful"
)

func init() {
	db := db.InitDB(models.DataSourceName(config.GetConfig().DATABASE_URL))
	mqConn := mq.InitMQ(models.MQUrl(config.GetConfig().MQ_URL))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit

		log.Println("Gracefully shutting down...")

		db.Close()
		mqConn.Close()

		os.Exit(0)
	}()
}

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
