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

func GetDatabase() (*mongo.Database, error) {
	mongoOnce.Do(initMongoDb)

	return mongoDatabase, mongoError
}

func initMongoDb() {
	values := url.Values{}
	values.Set("retryWrites", "true")
	values.Set("w", "majority")

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s",
		config.MongoUsr,
		config.MongoPwd,
		config.MongoUri,
		config.MongoDb,
		values.Encode(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
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
