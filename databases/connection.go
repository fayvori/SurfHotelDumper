package databases

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var Ctx = context.TODO()

var (
	Client *mongo.Client
)

func SetupMongoClient() {
	const uri = "mongodb://localhost:27017"

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Printf(err.Error())
	}

	defer func() {
		if err := client.Disconnect(Ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(Ctx, readpref.Primary()); err != nil {
		log.Fatalf(err.Error())
	}

	Client = client

	fmt.Println("Successfully connected to the database.")
}
