package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"dahbura.me/api/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDatabase *mongo.Database
	mongoError    error
	mongoOnce     sync.Once
)

func GetMongoDb() (*mongo.Database, error) {
	mongoOnce.Do(initMongoDb)

	return mongoDatabase, mongoError
}

func initMongoDb() {
	opts := url.Values{}
	opts.Set("retryWrites", "true")
	opts.Set("w", "majority")

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s",
		config.MongoUsr,
		config.MongoPwd,
		config.MongoUri,
		config.MongoDb,
		opts.Encode(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		mongoDatabase = nil
		mongoError = err
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		mongoDatabase = nil
		mongoError = err
		return
	}

	mongoDatabase = client.Database(config.MongoDb)
	mongoError = nil
}
