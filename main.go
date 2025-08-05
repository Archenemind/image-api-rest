package main

import (
	"github.com/Archenemind/image-api-rest/routers"

	_ "github.com/gin-gonic/gin"
)

func main() {
	router := routers.SetupRouters()

	router.Run("localhost:8080")

}
