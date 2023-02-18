package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

func New() *Server {
	engine := gin.Default()
	engine.Delims("{[{", "}]}")

	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return &Server{engine}
}

func (s *Server) Run(addr string) error {
	return s.Engine.Run(addr)
}

func (s *Server) RunTLS(addr string, certFile string, keyFile string) error {
	return s.Engine.RunTLS(addr, certFile, keyFile)
}
