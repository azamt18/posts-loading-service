package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"postsloading/pkg/models"
)

var response []models.Post

type PostsController interface {
	LoadPosts(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	Db *mongo.Collection
}

func (c controller) LoadPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Load posts request...")

	// make an API request to load posts
	params := url.Values{}
	params.Add("page", "1")

	resp, err := http.Get("https://gorest.co.in/public/v2/posts?" + params.Encode())
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)

	// Unmarshal result
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

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
		res, err := c.Db.InsertOne(context.TODO(), postObject)
		if err != nil {
			log.Fatalf("Internal error: %v", err)
		}

		// Check the insertion result
		objectId, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Fatalf("Can not convert to ObjectId")
		}

		log.Printf("Inserted objectId: %v", objectId)
	}

}

func NewController(db *mongo.Collection) PostsController {
	return controller{db}
}
