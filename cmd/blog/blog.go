package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kurnosovmak/short-link/internal/pkg/config"
	httpserver "github.com/kurnosovmak/short-link/pkg/http-server"
	"github.com/kurnosovmak/short-link/pkg/logging"
)

func main() {
	logging.Init()
	log := logging.GetLogger()
	log.Println("logger initialized")

	log.Println("config initializing")
	cfg := config.Get()

	log.Println("http server initializing")
	httpServer := httpserver.New()

	httpServer.GET("/api/heals", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Println("start server")
	httpServer.Run(cfg.GetFullAddress())

}
