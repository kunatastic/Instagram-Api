package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func create_new_user(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		response.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(request.Body).Decode(&user)
		collection := client.Database("Instagram").Collection("Users")
		user.Password = fmt.Sprint(md5.Sum([]byte(user.Password)))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(response).Encode(result)

	default:
		errorResponse(response, http.StatusMethodNotAllowed, "Invalid Request Method")
	}
}

func get_existing_user(response http.ResponseWriter, request *http.Request) {
	// fmt.Println(r.Method)
	switch request.Method {
	case "GET":
		response.Header().Add("content-type", "application/json")
		URL := request.URL.Path
		params := strings.TrimLeft(URL, "/users/")
		id, _ := primitive.ObjectIDFromHex(params)
		var user User
		collection := client.Database("Instagram").Collection("Users")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(user)

	default:
		errorResponse(response, http.StatusMethodNotAllowed, "Invalid Request Method")
	}
}

func create_new_post(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		response.Header().Add("content-type", "application/json")
		var post Post
		json.NewDecoder(request.Body).Decode(&post)
		collection := client.Database("Instagram").Collection("Posts")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		result, _ := collection.InsertOne(ctx, post)
		json.NewEncoder(response).Encode(result)

	default:
		errorResponse(response, http.StatusMethodNotAllowed, "Invalid Request Method")
	}
}

func get_existing_post(response http.ResponseWriter, request *http.Request) {
	// fmt.Println(r.Method)
	switch request.Method {
	case "GET":
		response.Header().Add("content-type", "application/json")
		URL := request.URL.Path
		params := strings.TrimLeft(URL, "/posts/")
		id, _ := primitive.ObjectIDFromHex(params)
		var post Post
		collection := client.Database("Instagram").Collection("Posts")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(post)

	default:
		errorResponse(response, http.StatusMethodNotAllowed, "Invalid Request Method")
	}
}

func posts_of_user(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		response.Header().Add("content-type", "application/json")
		URL := request.URL.Path
		params := strings.TrimLeft(URL, "/posts/users/")
		id, _ := primitive.ObjectIDFromHex(params)
		items, err := strconv.Atoi(params)
		var posts []Post
		collection := client.Database("Instagram").Collection("Posts")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := collection.Find(ctx, Post{UserID: id})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		var i int
		defer cursor.Close(ctx)
		for cursor.Next(ctx) && i < items {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
			i++
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(posts)

	default:
		errorResponse(response, http.StatusMethodNotAllowed, "Invalid Request Method")
	}
}
