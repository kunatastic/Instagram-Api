package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL  string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Timestamp string             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

var client *mongo.Client

func Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Init() {
	var mongoURI string = "mongodb://localhost:27017/instagram"

	// setup connection with mongoDb cluster
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	fmt.Println("Connected to MongoDB!")

	createRoutes()
}

func createRoutes() {
	http.HandleFunc("/users", create_new_user)
	http.HandleFunc("/users/", get_existing_user)
	http.HandleFunc("/posts", create_new_post)
	http.HandleFunc("/posts/", get_existing_post)
	http.HandleFunc("/posts/users/", posts_of_user)
}
