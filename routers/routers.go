package routers

import "github.com/gin-gonic/gin"

func SetupRouters() *gin.Engine {
	router := gin.Default()

	SetupImageRoutes(router)

	return router
}
