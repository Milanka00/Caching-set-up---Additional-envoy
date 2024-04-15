package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var payloads map[int][]byte
var once sync.Once
var collection *mongo.Collection

type Person struct {
	ID   int    `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

func main() {
	generatePayloads()

	// MongoDB connection
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	collection = client.Database("test").Collection("people")

	router := mux.NewRouter()

	router.HandleFunc("/publiccache/{id}", PublicCacheHandler).Methods("GET")
	router.HandleFunc("/publiccache/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/publiccache/{id}", UpdatePerson).Methods("PUT")
	router.HandleFunc("/publiccache/{id}", DeletePerson).Methods("DELETE")

	fmt.Println("Server is listening on port 8083...")
	http.ListenAndServe(":8083", router)
}

func PublicCacheHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payload, ok := payloads[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=600")
	w.Write(payload)
	additionalContent := []byte(" cached as public for ID " + strconv.Itoa(id))
	w.Write(additionalContent)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = collection.InsertOne(context.Background(), person)
	if err != nil {
		http.Error(w, "Error creating person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedPerson Person
	err = json.NewDecoder(r.Body).Decode(&updatedPerson)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedPerson}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error deleting person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generatePayloads() {
	once.Do(func() {
		payloads = make(map[int][]byte)
		for i := 1; i <= 100; i++ {
			payload := make([]byte, 1000)
			for j := range payload {
				payload[j] = 'x'
			}
			payloads[i] = payload
		}
	})
}
