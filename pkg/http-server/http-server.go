package httpserver

import (
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*gin.Engine
}

func New() *HttpServer {
	server := HttpServer{
		gin.Default(),
	}
	return &server
}
