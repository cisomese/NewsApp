package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newsapp/news/config"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type News struct {
	ID       int64     `json:"id" bson:"_id"`
	Date     time.Time `json:"date" bson:"date"`
	Source   string    `json:"source" bson:"source"`
	Message  string    `json:"message" bson:"message"`
	Path     string    `json:"path" bson:"media_path"`
	Category string    `json:"category" bson:"category"`
}

var collection *mongo.Collection

func init() {
	client, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	collection = client.Database(os.Getenv("DB_NAME")).Collection("news")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check for ?id=123 query parameter
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		// Fetch single news by ID
		newsID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var news News
		err = collection.FindOne(ctx, bson.M{"_id": newsID}).Decode(&news)
		if err == mongo.ErrNoDocuments {
			http.Error(w, fmt.Sprintf("no news found with id %d", newsID), http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(news)
		return
	}

	// If no id, fetch all news
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var newsList []News
	if err := cursor.All(ctx, &newsList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsList)
}
