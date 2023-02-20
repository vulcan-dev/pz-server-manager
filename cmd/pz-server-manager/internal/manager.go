package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vulcan-dev/pz-server-manager/cmd/pz-server-manager/internal/pz"
	"path/filepath"
	"strings"
)

type Config struct {
	HTTP_PORT string
	PZ_ROOT   string
}

var config Config
var options pz.Options
var fs pz.Filesystem

func Run(_config Config) error {
	config = _config

	log.Info("Starting HTTP server on port ", config.HTTP_PORT)
	log.Info("PZ_ROOT is ", config.PZ_ROOT)

	var err error
	options, err = pz.Parse("servertest", config.PZ_ROOT)
	if err != nil {
		log.Fatal(err)
	}

	fs = *pz.NewFilesystem(config.PZ_ROOT)
	_, err = fs.List("", false) // Load cache
	if err != nil {
		log.Fatal(err)
	}

	StartHTTPServer()
	return nil
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

	// General
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Filesystem
	engine.GET("/api/v1/files", func(c *gin.Context) {
		files, err := fs.List("", false)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   files,
		})
	})

	// Example: /api/v1/files/levels/level1
	engine.GET("/api/v1/files/*path", func(c *gin.Context) {
		path := c.Param("path")

		// Check if path is a file
		if filepath.Ext(strings.TrimSpace(path)) != "" {
			file, err := fs.Get(path)
			if err != nil {
				c.JSON(404, gin.H{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}

			content, err := file.Content()
			if err != nil {
				c.JSON(404, gin.H{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"status": "ok",
				"data":   content,
			})
			return
		}

		files, err := fs.List(path, false)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   files,
		})
	})

	// Config
	engine.GET("/api/v1/config", func(c *gin.Context) {
		data, err := options.JSON()
		if err != nil {
			log.Error(err)
			c.JSON(500, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   data,
		})
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

		log.Info(data)
	})

	engine.Run(":" + config.HTTP_PORT)
}
