package driver

import (
	"context"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models/datastore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
	"time"
)

type MongoDriver struct {
	Client         *mongo.Client
	collection     *mongo.Collection
	dbName         string
	collectionName string
}

// New function returns an instance of Engine and sets database connection
func New(uri, dbName, collectionName string) *MongoDriver {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed while initializing db client", zap.Error(err))
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed while connecting to db" + err.Error())
	}

	// TODO : Defer in main.go
	//defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping the db", zap.Error(err))
	}

	// Access a MongoDB collection through a database
	col := client.Database(dbName).Collection(collectionName)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	// Creates index of type `text`
	models := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "title", Value: bsonx.String("text")},
				{Key: "desc", Value: bsonx.String("text")}},
		},
		{
			Keys: bsonx.Doc{{Key: "publishedAt", Value: bsonx.Int32(-1)}},
		},
	}

	_, err = col.Indexes().CreateMany(ctx, models, opts)
	// Check for the options errors
	if err != nil {
		log.Fatal("Failed while creating indexes", zap.Error(err))
	}

	log.Info("Created indexes", zap.Any("Options", opts))
	return &MongoDriver{
		Client:         client,
		collection:     col,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (md MongoDriver) GetAll(limit, page int) (int64, []datastore.YTRecord, error) {
	result := make([]datastore.YTRecord, 0)
	ctx := context.Background()

	count, err := md.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Error("Error while counting documents from db.", zap.Error(err))
		return count, result, err
	}
	log.Info("Number of records returned", zap.Int64("Count", count))

	opts := NewMongoPaginate(limit, page).GetPaginatedOpts()
	opts.SetSort(bson.D{{"publishedAt", -1}})

	curr, err := md.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Error("Error while finding documents from db.", zap.Error(err))
		return count, result, err
	}

	for curr.Next(ctx) {
		var record datastore.YTRecord
		if err := curr.Decode(&record); err != nil {
			log.Error("Failed while decoding response from db", zap.Error(err))
		}

		result = append(result, record)
	}

	return count, result, nil
}

func (md MongoDriver) GetByTitleOrDesc(limit, page int, search string) (int64, []datastore.YTRecord, error) {
	result := make([]datastore.YTRecord, 0)
	ctx := context.Background()

	filter := bson.D{{"$text", bson.D{{"$search", search}}}}

	count, err := md.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Error("Error while counting documents from db.", zap.Error(err))
		return count, result, err
	}
	log.Info("Number of records returned", zap.Int64("Count", count))

	opts := NewMongoPaginate(limit, page).GetPaginatedOpts()
	opts.SetSort(bson.D{{"publishedAt", -1}})

	curr, err := md.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Error("Error while finding documents from db.", zap.Error(err))
		return count, result, err
	}

	for curr.Next(ctx) {
		var record datastore.YTRecord
		if err := curr.Decode(&record); err != nil {
			log.Info(err.Error())
		}

		result = append(result, record)
	}

	return count, result, nil
}

func (md MongoDriver) GetLastRecordTime() (*time.Time, error) {
	ctx := context.Background()

	options := options.Find()

	// Sort by `_id` field descending
	options.SetSort(bson.D{{"publishedAt", -1}})
	options.SetLimit(1)

	curr, err := md.collection.Find(ctx, bson.D{}, options)
	if err != nil {
		log.Error("Error while finding documents from db.", zap.Error(err))
		return nil, err
	}

	var lastRecordTime *time.Time

	for curr.Next(ctx) {
		var record datastore.YTRecord
		if err := curr.Decode(&record); err != nil {
			log.Info(err.Error())
		}

		lastRecordTime = &record.PublishedAt
	}

	return lastRecordTime, nil
}

func (md MongoDriver) SaveAll(records []datastore.YTRecord) error {
	ctx := context.Background()

	recordsToBeInserted := make([]interface{}, len(records))
	for i, v := range records {
		recordsToBeInserted[i] = v
	}

	_, err := md.collection.InsertMany(ctx, recordsToBeInserted)
	if err != nil {
		log.Error("Error while inserting documents into db.", zap.Error(err))
		return err
	}

	return nil
}
