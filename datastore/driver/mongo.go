package driver

import (
	"context"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models"
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

//func (md MongoDriver) Create() (*entities.Engine, error) {
//	_, err := es.DB.Exec("INSERT INTO engines VALUES (?,?,?,?)", e.ID, e.Displacement, e.Noc, e.Range)
//
//	if err != nil {
//		return nil, &customerrors.Err{Msg: customerrors.EngineInsertionError}
//	}
//
//	return e, nil
//}

func (md MongoDriver) GetAll(limit, page int) (int64, []models.YTRecord, error) {
	result := make([]models.YTRecord, 0)
	ctx := context.Background()

	count, err := md.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return count, result, err
	}

	log.Info("Number of records returned", zap.Int64("Count", count))

	curr, err := md.collection.Find(ctx, bson.D{}, NewMongoPaginate(limit, page).GetPaginatedOpts())
	if err != nil {
		return count, result, err
	}

	for curr.Next(ctx) {
		var record models.YTRecord
		if err := curr.Decode(&record); err != nil {
			log.Info("Failed while decoding response from db", zap.Error(err))
		}

		result = append(result, record)
	}

	return count, result, nil
}

func (md MongoDriver) GetByTitleOrDesc(limit, page int, search string) (int64, []models.YTRecord, error) {
	result := make([]models.YTRecord, 0)
	ctx := context.Background()

	count, err := md.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return count, result, err
	}

	log.Info("Number of records returned", zap.Int64("Count", count))

	filter := bson.D{{"$text", bson.D{{"$search", search}}}}
	curr, err := md.collection.Find(ctx, filter, NewMongoPaginate(limit, page).GetPaginatedOpts())
	if err != nil {
		return count, result, err
	}
	for curr.Next(ctx) {
		var record models.YTRecord
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
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)

	curr, err := md.collection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}
	var lastRecordTime *time.Time

	for curr.Next(ctx) {
		var record models.YTRecord
		if err := curr.Decode(&record); err != nil {
			log.Info(err.Error())
		}

		lastRecordTime = &record.PublishedAt
	}

	return lastRecordTime, nil
}

func (md MongoDriver) SaveAll(records []models.YTRecord) error {
	ctx := context.Background()

	recordsToBeInserted := make([]interface{}, len(records))
	for i, v := range records {
		recordsToBeInserted[i] = v
	}

	_, err := md.collection.InsertMany(ctx, recordsToBeInserted)
	if err != nil {
		return err
	}

	return nil
}
