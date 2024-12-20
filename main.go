package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Medication defines the data structure for a medication item
// using struct tags for JSON marshaling/unmarshaling
type Medication struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

// getMedications is an http.HandlerFunc that returns all medications
func getMedications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch medications", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var medications []Medication
	if err = cursor.All(ctx, &medications); err != nil {
		http.Error(w, "Failed to parse medications", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(medications)
}

// getMedicationByID is an http.HandlerFunc that looks up a medication by ID
func getMedicationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var medication Medication
	err = db.FindOne(ctx, bson.M{"_id": id}).Decode(&medication)
	if err != nil {
		http.Error(w, "Medication not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(medication)
}

// createMedication is an http.HandlerFunc that adds a new medication
func createMedication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMedication Medication
	// Use json.Decoder to parse the request body directly into a struct
	if err := json.NewDecoder(r.Body).Decode(&newMedication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newMedication.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.InsertOne(ctx, newMedication)
	if err != nil {
		http.Error(w, "Failed to create medication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMedication)
}

func main() {
	// Connect to the database
	ConnectDB()
	defer CloseDB()

	// Create a new router
	r := mux.NewRouter()

	// Register route handlers with HTTP methods
	// Each HandleFunc returns a Route instance for further configuration
	r.HandleFunc("/medications", getMedications).Methods("GET")
	r.HandleFunc("/medications/{id}", getMedicationByID).Methods("GET")
	r.HandleFunc("/medications", createMedication).Methods("POST")

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
}
