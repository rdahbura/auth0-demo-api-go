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
	mongoClient *mongo.Client
	mongoError  error
	mongoOnce   sync.Once
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		initMongoClient()
	})

	return mongoClient, mongoError
}

func initMongoClient() {
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
		mongoClient = nil
		mongoError = err
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		mongoClient = nil
		mongoError = err
		return
	}

	mongoClient = client
	mongoError = nil
}
