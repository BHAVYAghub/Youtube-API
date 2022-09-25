package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BHAVYAghub/Youtube-API/datastore"
	"github.com/BHAVYAghub/Youtube-API/datastore/driver"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	api "github.com/BHAVYAghub/Youtube-API/routes"
	"github.com/BHAVYAghub/Youtube-API/scheduler"
	"github.com/BHAVYAghub/Youtube-API/service"
	"github.com/BHAVYAghub/Youtube-API/service/external"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	wg sync.WaitGroup
)

func main() {
	log.InitClient("info", "")

	// Load env configs
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file found")
	}

	// Initialize DB configs
	db := initializeDBDriver()
	defer func(Client *mongo.Client, ctx context.Context) {
		err := Client.Disconnect(ctx)
		if err != nil {
			log.Error("Error while disconnecting DB.")
		}
	}(db.Client, context.Background())

	// Initialize internal services
	svc := initializeYTSvc(db)

	// Initialize router
	router := gin.Default()
	api.AttachRoutes(router, svc)

	scheduleCron(svc)

	err := router.Run()
	if err != nil {
		log.Fatal("Error while running server.")
	}
}

func scheduleCron(svc service.Service) {
	disableCron := os.Getenv("DISABLE_CRON")
	if disableCron == "" {
		log.Warn("config DISABLE_CRON not set. ", zap.String("default", "false"))
		disableCron = "false"
	}

	disableCronBool, err := strconv.ParseBool(disableCron)
	if err != nil {
		log.Warn("Invalid config DISABLE_CRON. ", zap.String("default", "false"))
	}

	if disableCronBool {
		return
	}

	cronInterval := os.Getenv("YT_API_FETCH_INTERVAL")
	if cronInterval == "" {
		log.Warn("config YT_API_FETCH_INTERVAL not set. ", zap.String("default", "1"))
		cronInterval = "1"
	}

	cronIntervalInt, err := strconv.Atoi(cronInterval)
	if err != nil {
		log.Warn("Invalid config YT_API_FETCH_INTERVAL. ", zap.String("default", "1"))
		cronIntervalInt = 1
	}

	//Call the ticker function before submitting it to the scheduler
	if err := svc.FetchAndInsertRecords(); err != nil {
		log.Warn("Failed while inserting data in db", zap.Error(err))
	}

	go scheduler.SubmitToTicker(&wg, svc.FetchAndInsertRecords, time.Duration(cronIntervalInt)*time.Minute)
}

func initializeDBDriver() *driver.MongoDriver {
	missingConfigs := make([]string, 0)

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		missingConfigs = append(missingConfigs, "MONGODB_URI")
	}

	dbName := os.Getenv("MONGODB_DATABASE_NAME")
	if dbName == "" {
		missingConfigs = append(missingConfigs, "MONGODB_DATABASE_NAME")
	}

	colName := os.Getenv("MONGODB_COLLECTION_NAME")
	if colName == "" {
		missingConfigs = append(missingConfigs, "MONGODB_COLLECTION_NAME")
	}

	if len(missingConfigs) > 0 {
		log.Fatal("Missing mandatory configs. ", zap.Any("configs", missingConfigs))
	}
	return driver.New(mongoURI, dbName, colName)
}

func initializeYTSvc(driver datastore.Database) service.Service {
	missingConfigs := make([]string, 0)

	ytBaseUrl := os.Getenv("YT_API_BASE_URL")
	if ytBaseUrl == "" {
		missingConfigs = append(missingConfigs, "YT_API_BASE_URL")
	}

	ytQueryString := os.Getenv("YT_QUERY_STRING")
	if ytQueryString == "" {
		missingConfigs = append(missingConfigs, "YT_QUERY_STRING")
	}

	ytAPIKey := os.Getenv("YT_API_KEY")
	if ytAPIKey == "" {
		missingConfigs = append(missingConfigs, "YT_API_KEY")
	}

	fetchRecordsAfter := os.Getenv("YT_FETCH_RECORDS_AFTER")
	if fetchRecordsAfter == "" {
		missingConfigs = append(missingConfigs, "YT_FETCH_RECORDS_AFTER")
	}

	if len(missingConfigs) > 0 {
		log.Fatal("Missing mandatory configs. ", zap.Any("configs", missingConfigs))
	}

	fetchRecordsAfterTime, err := time.Parse(time.RFC3339, fetchRecordsAfter)
	if err != nil {
		log.Fatal("Invalid time passed in FETCH_RECORDS_AFTER config.")
	}

	keys := strings.Split(ytAPIKey, ",")

	// Initializing YT client
	ytClient := external.NewService(ytBaseUrl, ytQueryString, keys)

	return service.NewYTService(ytClient, fetchRecordsAfterTime, driver)
}
