package http

import "github.com/gin-gonic/gin"

// new gin router
func NewRouter(imgService ImageService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	imgController := newImageController(imgService)
	api := router.Group("/api")
	{
		api.POST("/upload", imgController.upload)
		api.GET("/download/:id", imgController.download)
		api.DELETE("/delete/:id", imgController.delete)
	}
	return router
}
