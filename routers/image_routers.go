package routers

import (
	"github.com/Archenemind/image-api-rest/controllers"

	"github.com/gin-gonic/gin"
)

func SetupImageRoutes(router *gin.Engine) {
	images := router.Group("/images")
	{
		images.GET("/images", controllers.GetImages)
		images.GET("/image/:id", controllers.GetImageById)
		images.POST("/upload", controllers.UploadImage)
		images.POST("/convert", controllers.ConvertAndDeleteImage)
		images.PUT("/image/:id", controllers.UpdateImage)
		images.DELETE("/image/:id", controllers.DeleteImage)
	}
}
