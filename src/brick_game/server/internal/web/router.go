package web

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// NewRouter creates and returns a configured Gin router with API routes and static file serving.
func NewRouter(gh *GameHandler) *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/games", gh.GetGames)
		api.POST("/games/:id", gh.StartGame)
		api.POST("/actions", gh.PostAction)
		api.GET("/state", gh.GetState)
	}

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	router.StaticFS("/static", gin.Dir(filepath.Join(exeDir, "../brick_game/gui/web_gui"), false))
	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(exeDir, "../brick_game/gui/web_gui/index.html"))
	})

	return router
}
