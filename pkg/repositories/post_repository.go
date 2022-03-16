package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"postsloading/pkg/db"
	"postsloading/pkg/models"
)

type PostRepository struct {
	client *mongo.Client
}

func (repository *PostRepository) getClient() *mongo.Client {
	if repository.client != nil {
		return repository.client
	}

	repository.client = db.ConnectToMongo()
	return repository.client
}

func (repository *PostRepository) SaveToDb(response []models.Post) (bool, error) {
	collection := repository.getClient().Database("mydb").Collection("posts")

	// Iterate over the values
	for i, post := range response {
		fmt.Printf("%v) Post title: %v\n", i, post.Title)
		postObject := models.Post{
			Id:     post.Id,
			UserId: post.UserId,
			Title:  post.Title,
			Body:   post.Body,
		}

		// Saving to db
		res, err := collection.InsertOne(context.TODO(), postObject)
		if err != nil {
			log.Fatalf("Internal error: %v", err)
			return false, err
		}

		// Check the insertion result
		objectId, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Fatalf("Can not convert to ObjectId")
			return false, err
		}

		log.Printf("Inserted objectId: %v", objectId)
	}

	return true, nil
}
