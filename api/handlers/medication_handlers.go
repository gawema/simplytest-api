package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"simplytest-api/storage/models"

	"github.com/gorilla/mux"
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (h *MedicationHandler) UpdateMedication(w http.ResponseWriter, r *http.Request) {
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

func (h *MedicationHandler) DeleteMedication(w http.ResponseWriter, r *http.Request) {
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
