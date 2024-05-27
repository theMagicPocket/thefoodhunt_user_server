package dbcontext

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDB(uri string) (*mongo.Database, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		return nil, err
	}

	mongoClient = client
	return mongoClient.Database("yumfoods"), nil

}

func DisconnectDB() error {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
