package routers

import (
	"github.com/Archenemind/image-api-rest/controllers"
	"github.com/Archenemind/image-api-rest/utils"
	"github.com/gin-gonic/gin"
)

func SetupUsersRouters(router *gin.Engine) {
	users := router.Group("/users")
	{
		// users.GET("/", controllers.GetUsers)
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.LoginUser)
		users.POST("/refresh", controllers.RefreshToken)
		users.PUT("/update/:id", utils.JWTMiddleware(), controllers.UpdateUser)
		// users.GET("/:id", controllers.GetUser)
		// users.PUT("/:id", controllers.UpdateUser)
		// users.DELETE("/:id", controllers.DeleteUser)
	}
}
