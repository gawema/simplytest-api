package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"simplytest-api/storage/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MedicationHandler handles HTTP requests for medication operations
type MedicationHandler struct {
	collection *mongo.Collection // MongoDB collection for medications
}

// NewMedicationHandler creates a new handler with the given MongoDB collection
func NewMedicationHandler(collection *mongo.Collection) *MedicationHandler {
	return &MedicationHandler{
		collection: collection,
	}
}

// enableCORS adds CORS headers to the response
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // For development. In production, set to your Next.js domain
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// GetMedications handles GET requests to fetch all medications
func (h *MedicationHandler) GetMedications(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all medications in the database
	cursor, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch medications", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of medications
	var medications []models.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		http.Error(w, "Failed to parse medications", http.StatusInternalServerError)
		return
	}

	// Return the medications as JSON
	json.NewEncoder(w).Encode(medications)
}

// GetMedicationByID handles GET requests to fetch a single medication by ID
func (h *MedicationHandler) GetMedicationByID(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL parameters
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the medication in the database
	var medication models.Medication
	err = h.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&medication)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Medication not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch medication", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(medication)
}

// CreateMedication handles POST requests to create a new medication
func (h *MedicationHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a Medication struct
	var newMedication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&newMedication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Let's add some debug logging
	fmt.Printf("Received medication: %+v\n", newMedication)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the new medication into the database
	result, err := h.collection.InsertOne(ctx, newMedication)
	if err != nil {
		http.Error(w, "Failed to create medication", http.StatusInternalServerError)
		return
	}

	// Set the ID of the newly created medication
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newMedication.ID = oid
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMedication)
}

// UpdateMedication handles PUT requests to update a medication
func (h *MedicationHandler) UpdateMedication(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	var medication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        medication.Name,
			"description": medication.Description,
			"price":       medication.Price,
			"imageUrl":    medication.ImageURL,
		},
	}

	result, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update medication", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Medication not found", http.StatusNotFound)
		return
	}

	medication.ID = id
	json.NewEncoder(w).Encode(medication)
}

// DeleteMedication handles DELETE requests to delete a medication
func (h *MedicationHandler) DeleteMedication(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete medication", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Medication not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
