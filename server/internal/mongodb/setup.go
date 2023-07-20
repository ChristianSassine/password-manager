package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoUsername = errors.New("missing MONGO_USERNAME environmental variable for MongoDB")
	ErrMongoPassword = errors.New("missing MONGO_PASSWORD environmental variable for MongoDB")
	ErrMongoAddr     = errors.New("missing MONGO_ADDR environmental variable for MongoDB")
	ErrMongoPort     = errors.New("missing MONGO_PORT environmental variable for MongoDB")
)

var (
	client     *mongo.Client
	clientDb   *mongo.Database
	collection *mongo.Collection
)

const (
	databaseName   = "password-manager"
	userCollection = "pass-manager"
)

type userData struct {
	Username  string
	Password  string
	passwords map[string]string
}

func Start() error {
	username, ok := os.LookupEnv("MONGO_USERNAME")
	if !ok {
		log.Fatal(ErrMongoUsername)
	}
	password, ok := os.LookupEnv("MONGO_PASSWORD")
	if !ok {
		log.Fatal(ErrMongoPassword)
	}
	ip, ok := os.LookupEnv("MONGO_ADDR")
	if !ok {
		log.Fatal(ErrMongoAddr)
	}
	port, ok := os.LookupEnv("MONGO_PORT")
	if !ok {
		log.Fatal(ErrMongoPort)
	}

	var uri = fmt.Sprintf("mongodb://%v:%v@%v:%v/", username, password, ip, port)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
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
	log.Println("Successfully connected to MongoDB!")

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
