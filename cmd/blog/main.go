package main

import (
	"context"
	"github.com/kurnosovmak/short-link/internal/handlers/links"
	"github.com/kurnosovmak/short-link/internal/pkg/config"
	"github.com/kurnosovmak/short-link/internal/services/link_service"
	"github.com/kurnosovmak/short-link/pkg/cachemap"
	httpserver "github.com/kurnosovmak/short-link/pkg/http-server"
	"github.com/kurnosovmak/short-link/pkg/logging"
	"github.com/kurnosovmak/short-link/pkg/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ctx = context.Background()

// @title           Shorter link backend
// @version         1.0
// @description     This is a sample service shorter link.

// @contact.name   Kurnosovmak
// @contact.email  kurnosovmak@gmail.com

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api/

func main() {
	logging.Init()
	log := logging.GetLogger()
	log.Println("logger initialized")

	log.Println("config initializing")
	cfg := config.Get()

	log.Println("redis connect initializing")
	redis := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	log.Println("cahce initializing")
	cacheLink := cachemap.New(log)

	log.Println("http server initializing")
	server := httpserver.New()

	log.Println("create and register handlers")
	linkService := link_service.NewService(redis, cacheLink, log)
	linkHandler := links.NewHandler(linkService, log)
	linkHandler.Register(server)

	if *cfg.IsDebug {
		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Println("start server")
	server.Run(cfg.GetFullAddress())

}
