package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	URI = "mongodb://localhost:27017"
)

var client *mongo.Client

// In Golang, init() functions always initialize whenever the package is called.
// So, whenever DB variable called, the init() function initialized
func init() {
	ConnectToMongo()
}

func ConnectToMongo() *mongo.Client {
	fmt.Println("Connection to MongoDb")
	// Create client
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}

	// Create connect
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseClientDB() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
}
