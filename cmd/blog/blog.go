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
)

var ctx = context.Background()

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

	log.Println("start server")
	server.Run(cfg.GetFullAddress())

}
