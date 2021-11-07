package mongodb

import (
	"context"
	"fmt"
	"time"

	errors "github.com/mongodb-developer/docker-golang-example/Errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDbClient() (client *mongo.Client, ctx context.Context) {
	fmt.Println("Starting the application")

	client, err := mongo.NewClient(options.Client().
		ApplyURI("******************************************")) //add mongo uri
	if err != nil {
		errors.HandleCustomError(errors.FetchClientError, err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		errors.HandleCustomError(errors.ConnectingDbError, err)
	}
	return client, ctx
}
