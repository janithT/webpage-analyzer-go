package engine

import (
	"log"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	httpHandler "github.com/janithT/webpage-analyzer/handler/http"
	"github.com/janithT/webpage-analyzer/handler/middleware"
)

// NewRouter sets up Gin routes and returns the router instance
func NewRouter() *gin.Engine {
	log.Println("entry point - Gin router - j1")

	router := gin.Default()

	// CORS
	router.Use(middleware.SetupCORS())

	// Serve Angular dist output
	// router.Static("/", "./web/wep-page-analyzer-ng")

	// API route - use api prefix later
	router.GET("/v1/analyze", httpHandler.AnalyzeHandler)

	// fallback angular
	router.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./web/index.html")
		} else {
			c.File("./web/" + path.Join(dir, file))
		}
	})

	return router
}
