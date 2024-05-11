package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func SetupMongoDB(uri string, db string) (*mongo.Client, *mongo.Database, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	database := client.Database(db)

	return client, database, ctx, cancel, nil
}

// Close the connection
func CloseConnection(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer func() error {
		cancel()
		if err := client.Disconnect(context); err != nil {
			return err
		}
		log.Println("Close connection is called")
		return nil
	}()

}
