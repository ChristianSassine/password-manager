package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Command wrappers for other modules
func Add(document interface{}, opts ...*options.InsertOneOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	passCollection.InsertOne(ctx, document, opts...)
}

func Remove(filter interface{}, opts ...*options.DeleteOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	passCollection.DeleteOne(ctx, filter, opts...)
}

func Modify(filter interface{}, update interface{}, opts ...*options.UpdateOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	passCollection.UpdateOne(ctx, filter, update, opts...)
}

func Get(filter interface{}, opts ...*options.FindOneOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	passCollection.FindOne(ctx, filter, opts...)
}
