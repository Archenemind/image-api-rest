package main

import (
	"github.com/Archenemind/api-rest/routers"

	_ "github.com/gin-gonic/gin"
)

// images := make([]image,5)

func main() {
	router := routers.SetupRouters()

	router.Run("localhost:8080")

}
