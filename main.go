package main

import (
	log "github.com/BHAVYAghub/Youtube-API/logging"
	api "github.com/BHAVYAghub/Youtube-API/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api.AttachRoutes(router)

	err := router.Run()
	if err != nil {
		log.Fatal("Error while running server")
	}
}
