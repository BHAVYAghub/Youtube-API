package driver

import (
	"context"
	"fmt"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"os"
	"time"
)

type MongoDriver struct {
	client *mongo.Client
}

// New function returns an instance of Engine and sets database connection
func New(atlasURI string) MongoDriver {

	client, err := mongo.NewClient(options.Client().ApplyURI(atlasURI))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err.Error())
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(databases)
	fmt.Println("thsis")

	models := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "title", Value: bsonx.String("text")},
				{Key: "desc", Value: bsonx.String("text")}},
		},
	}

	// Access a MongoDB collection through a database
	col := client.Database("test1").Collection("text-test")
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err = col.Indexes().CreateMany(ctx, models, opts)

	// Check for the options errors
	if err != nil {
		fmt.Println("Indexes().CreateIndexes() ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		fmt.Println("CreateIndexes() opts:", opts)
	}

	return MongoDriver{client: client}
}
