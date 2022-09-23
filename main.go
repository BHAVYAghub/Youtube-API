package main

import (
	"github.com/BHAVYAghub/Youtube-API/datastore"
	"github.com/BHAVYAghub/Youtube-API/datastore/driver"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	api "github.com/BHAVYAghub/Youtube-API/routes"
	"github.com/BHAVYAghub/Youtube-API/service"
	"github.com/BHAVYAghub/Youtube-API/service/external"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file found")
	}

	driver := initializeDBDriver()

	svc := initializeYTSvc(driver)

	router := gin.Default()
	api.AttachRoutes(router, svc)

	err := router.Run()
	if err != nil {
		log.Fatal("Error while running server.")
	}
}

func initializeDBDriver() *driver.MongoDriver {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE_NAME")
	colName := os.Getenv("MONGODB_COLLECTION_NAME")

	return driver.New(mongoURI, dbName, colName)
}

func initializeYTSvc(driver datastore.Database) service.Service {
	ytBaseUrl := os.Getenv("YT_API_BASE_URL")
	ytQueryString := os.Getenv("YT_QUERY_STRING")
	ytAPIKey := os.Getenv("YT_API_KEY")
	ytClient := external.NewService(ytBaseUrl, ytQueryString, ytAPIKey)

	fetchRecordsAfter, err := time.Parse(time.RFC3339, os.Getenv("FETCH_RECORDS_AFTER"))
	if err != nil {
		log.Fatal("Invalid time passed in FETCH_RECORDS_AFTER config.")
	}

	return service.NewYTService(ytClient, fetchRecordsAfter, driver)
}
