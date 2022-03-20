package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId int                `json:"user_id,omitempty"`
	Title  string             `json:"title,omitempty"`
	Body   string             `json:"body,omitempty"`
}
