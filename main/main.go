package main

import (
	"api-rest/controllers"

	"github.com/gin-gonic/gin"
)

// images := make([]image,5)

func main() {
	router := gin.Default()
	router.GET("/images", controllers.GetImages)
	router.GET("/image/:id", controllers.GetImageById)
	router.POST("/image", controllers.PostImages)
	router.POST("/upload", controllers.UploadImage)
	router.POST("/convert", controllers.ConvertImage)
	router.PUT("/image/:id", controllers.UpdateImage)
	router.DELETE("/image/:id", controllers.DeleteImage)

	router.Run("localhost:8080")

}
