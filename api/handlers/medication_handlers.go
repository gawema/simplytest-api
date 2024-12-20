package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"simplytest-api/storage/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MedicationHandler struct {
	collection *mongo.Collection
}

func NewMedicationHandler(collection *mongo.Collection) *MedicationHandler {
	return &MedicationHandler{
		collection: collection,
	}
}

func (h *MedicationHandler) GetMedications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch medications", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var medications []models.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		http.Error(w, "Failed to parse medications", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(medications)
}

func (h *MedicationHandler) GetMedicationByID(w http.ResponseWriter, r *http.Request) {
	// ... similar to before but using h.collection ...
}

func (h *MedicationHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newMedication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&newMedication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.collection.InsertOne(ctx, newMedication)
	if err != nil {
		http.Error(w, "Failed to create medication", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newMedication.ID = oid
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMedication)
}
