package main

import (
	"github.com/BHAVYAghub/Youtube-API/driver"
	api "github.com/BHAVYAghub/Youtube-API/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	_ = driver.New(os.Getenv("MONGODB_URI"))

	router := gin.Default()
	api.AttachRoutes(router)

	err := router.Run()
	if err != nil {
		log.Fatal("Error while running server")
	}
}
