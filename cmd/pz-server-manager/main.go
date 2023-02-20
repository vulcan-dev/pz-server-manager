package main

import (
	log "github.com/sirupsen/logrus"
	server "github.com/vulcan-dev/pz-server-manager/cmd/pz-server-manager/internal"
	"os"
)

var config server.Config

func main() {
	// Init logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Load config
	config.HTTP_PORT = os.Getenv("HTTP_PORT")
	config.PZ_ROOT = os.Getenv("PZ_ROOT")
	if config.PZ_ROOT == "" {
		log.Fatal("PZ_ROOT not set")
	}

	// Start server
	server.Run(config)
}
