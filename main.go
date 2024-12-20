package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Medication defines the data structure for a medication item
// using struct tags for JSON marshaling/unmarshaling
type Medication struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// medications is an in-memory slice that stores Medication instances
// TODO: This is not thread-safe and would need mutex protection for concurrent access
var medications = []Medication{
	{ID: 1, Name: "Paracetamol", Description: "Pain reliever and fever reducer."},
	{ID: 2, Name: "Ibuprofen", Description: "Reduces inflammation and treats pain."},
}

// nextID is used for auto-incrementing medication IDs
// TODO: In production, this should be handled by a database
var nextID = 3

// getMedications is an http.HandlerFunc that returns all medications
// as a JSON array in the response body
func getMedications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medications)
}

// getMedicationByID is an http.HandlerFunc that looks up a medication by ID
// It extracts the ID from the URL path using gorilla/mux router variables
func getMedicationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// Convert string ID to int using strconv package
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	// Linear search through medications slice
	// TODO:In a real app, this would be a database query
	for _, medication := range medications {
		if medication.ID == id {
			json.NewEncoder(w).Encode(medication)
			return
		}
	}

	http.Error(w, "Medication not found", http.StatusNotFound)
}

// createMedication is an http.HandlerFunc that adds a new medication
// It expects a JSON payload in the request body matching the Medication struct
func createMedication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMedication Medication
	// Use json.Decoder to parse the request body directly into a struct
	if err := json.NewDecoder(r.Body).Decode(&newMedication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Append to medications slice and increment nextID
	newMedication.ID = nextID
	nextID++
	medications = append(medications, newMedication)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMedication)
}

func main() {
	// Create a new gorilla/mux router instance
	r := mux.NewRouter()

	// Register route handlers with HTTP methods
	// Each HandleFunc returns a Route instance for further configuration
	r.HandleFunc("/medications", getMedications).Methods("GET")
	r.HandleFunc("/medications/{id}", getMedicationByID).Methods("GET")
	r.HandleFunc("/medications", createMedication).Methods("POST")

	// Start HTTP server with the router as handler
	http.ListenAndServe(":8080", r)
}
