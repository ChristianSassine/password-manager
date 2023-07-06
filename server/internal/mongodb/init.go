package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientDb   *mongo.Database
	collection *mongo.Collection
)

const (
	databaseName   = "password-manager"
	userCollection = "DevData"
)

type userData struct {
	Username  string
	Password  string
	passwords map[string]string
}

func Start() error {
	var uri = "mongodb://root:pass@127.0.0.1:27017/" // TODO: Make the parameters environmental variables

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	var result bson.M
	clientDb = client.Database(databaseName)
	if err := clientDb.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MongoDB!") // TODO: Replace with actual logger

	if isCollectionMissing(userCollection) {
		createCollection()
	}

	collection = clientDb.Collection(userCollection)

	return nil
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func createCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result bson.M
	if err := clientDb.RunCommand(ctx, bson.D{{Key: "create", Value: userCollection}}).Decode(&result); err != nil {
		panic(err)
	}
}

func isCollectionMissing(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	colls, err := clientDb.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	for _, c := range colls {
		if c == name {
			return false
		}
	}
	return true
}
