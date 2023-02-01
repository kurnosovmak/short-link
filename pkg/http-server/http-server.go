package httpserver

import (
	"github.com/gin-gonic/gin"
)

type httpServer struct {
	*gin.Engine
}

func New() *httpServer {
	server := httpServer{
		gin.Default(),
	}
	return &server
}
