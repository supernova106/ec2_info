package main

import (
	"github.com/supernova106/ec2_info/app/config"
	"github.com/supernova106/ec2_info/app/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

var cfg *config.Config

func main() {
	// Load config
	var err error
	cfg, err = config.Load(".env")
	if err != nil {
		log.Fatalf("Can't load .env file %v", err)
		return
	}

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware

	router := gin.Default()
	router.Use(injectDependencyServices())

	router.GET("/", request.Check)

	router.GET("/price", request.GetData)

	// By default it serves on :8080 unless a
	// API_PORT environm+nt variable was defined.
	router.Run(":" + cfg.Port)
}

func injectDependencyServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	}
}
