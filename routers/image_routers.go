package routers

import (
	"github.com/Archenemind/image-api-rest/controllers"
	"github.com/Archenemind/image-api-rest/utils"
	"github.com/gin-gonic/gin"
)

func SetupImageRoutes(router *gin.Engine) {
	images := router.Group("/images")
	images.Use(utils.JWTMiddleware())
	{
		images.GET("/", controllers.GetImages)
		images.GET("/image/:id", controllers.GetImageById)
		images.POST("/upload", controllers.UploadImage)
		images.POST("/convert", controllers.ConvertAndDeleteImage)
		images.PUT("/image/:id", controllers.UpdateImage)
		images.DELETE("/image/:id", controllers.DeleteImage)
	}
}
