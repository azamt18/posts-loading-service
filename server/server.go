package server

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
	posts_loader "posts.loader/internal/protobuff/posts.loader.v1"
	"posts.loader/pkg/models"
)

type server struct {
	address string
	Db      *mongo.Collection
}

var response []models.Post

func (s server) LoadPosts(ctx context.Context, request *posts_loader.LoadPostsRequest) (*posts_loader.LoadPostsResponse, error) {
	fmt.Println("Load posts request...")

	// make an API request to load posts
	params := url.Values{}
	params.Add("page", "1")

	resp, err := http.Get("https://gorest.co.in/public/v2/posts?" + params.Encode())
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return nil, err
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
		return nil, err
	}

	// Iterate over the values
	loadedPostsCount := int64(0)
	for i, post := range response {
		fmt.Printf("%v) Post title: %v\n", i, post.Title)
		postObject := models.Post{
			Id:     post.Id,
			UserId: post.UserId,
			Title:  post.Title,
			Body:   post.Body,
		}

		// Saving to db
		res, err := s.Db.InsertOne(context.TODO(), postObject)
		if err != nil {
			log.Fatalf("Internal error: %v", err)
		}

		// Check the insertion result
		objectId, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Fatalf("Can not convert to ObjectId")
		}

		loadedPostsCount++
		log.Printf("Inserted objectId: %v", objectId)
	}

	return &posts_loader.LoadPostsResponse{LoadedPostsCount: loadedPostsCount}, nil
}

func NewGrpcServer(address string) posts_loader.PostsLoaderServiceServer {
	return &server{
		address: address,
	}
}
