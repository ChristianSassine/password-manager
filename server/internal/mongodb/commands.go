package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	maxCommandTime = 2 * time.Second
)

// Command wrappers for other packages
func Add(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxCommandTime)
	defer cancel()

	return collection.InsertOne(ctx, document, opts...)
}

func Remove(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxCommandTime)
	defer cancel()

	return collection.DeleteOne(ctx, filter, opts...)
}

func Update(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxCommandTime)
	defer cancel()

	return collection.UpdateOne(ctx, filter, update, opts...)
}

func Get(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), maxCommandTime)
	defer cancel()

	return collection.FindOne(ctx, filter, opts...)
}

func Exist(filter interface{}, opts ...*options.CountOptions) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxCommandTime)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
