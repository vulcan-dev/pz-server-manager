package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vulcan-dev/pz-server-manager/cmd/pz-server-manager/internal/pz"
)

type Config struct {
	HTTP_PORT string
	PZ_ROOT   string
}

var config Config
var pzConfig pz.Config

func Run(_config Config) {
	config = _config

	log.Info("Starting HTTP server on port ", config.HTTP_PORT)
	log.Info("PZ_ROOT is ", config.PZ_ROOT)

	var err error
	pzConfig, err = pz.Parse("servertest", config.PZ_ROOT)
	if err != nil {
		log.Fatal(err)
	}

	StartHTTPServer()
}

func StartHTTPServer() {
	engine := gin.Default()
	engine.Delims("{[{", "}]}")

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ! Important: Methods must be registered after middlewares

	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	engine.GET("/api/v1/config", func(c *gin.Context) {
		c.JSON(200, pzConfig.ReadToGIN(c))
	})

	engine.POST("/api/v1/config", func(c *gin.Context) {
		data := make(map[string]interface{})
		err := c.BindJSON(&data)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		err = pzConfig.SaveFromJSON("servertest", config.PZ_ROOT, data)
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		}
	})

	engine.Run(":" + config.HTTP_PORT)
}
