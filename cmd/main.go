package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"posts.loader/pkg/db"
	"posts.loader/pkg/handlers"
)

const (
	DBName          = "mydb"
	PostsCollection = "posts"
	URI             = "mongodb://localhost:27017"
)

var collection *mongo.Collection

func main() {
	// if the go code is crushed -> get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize connection to db
	collection = connectDb()

	// starting a service (for internal connections)
	fmt.Println("Start Posts Loading Service...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Panicf("FAILED TO LISTEN %v", err)
	}
	fmt.Printf("Posts Loading is listening on: %v", lis.Addr())

	// starting an API server (for public connections)
	if startApiServer(collection); err != nil {
		log.Panicf("FAILED TO START API SERVER: %v", err)
	}

	// closing the connection to db is not necessary
	// because of it is regulated
	defer db.CloseClientDB()

}

func connectDb() *mongo.Collection {
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

	collection = client.Database(DBName).Collection(PostsCollection)

	fmt.Println("Connected to database: ", collection.Database().Name())

	return collection
}

func startApiServer(collection *mongo.Collection) error {
	postsController := handlers.NewController(collection)
	router := mux.NewRouter()

	router.HandleFunc("/posts/load", postsController.LoadPosts).Methods(http.MethodGet)

	log.Println("API is running")
	return http.ListenAndServe(":4000", router)
}
