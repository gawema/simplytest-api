package api

import (
	"simplytest-api/api/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configures all the routes for our API
func SetupRoutes(r *mux.Router, collection *mongo.Collection) {
	// Create a new medication handler with the MongoDB collection
	medicationHandler := handlers.NewMedicationHandler(collection)

	// Add OPTIONS method to all routes
	r.HandleFunc("/medications", medicationHandler.GetMedications).Methods("GET", "OPTIONS")
	r.HandleFunc("/medications/{id}", medicationHandler.GetMedicationByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/medications", medicationHandler.CreateMedication).Methods("POST", "OPTIONS")
	r.HandleFunc("/medications/{id}", medicationHandler.UpdateMedication).Methods("PUT", "OPTIONS")
	r.HandleFunc("/medications/{id}", medicationHandler.DeleteMedication).Methods("DELETE", "OPTIONS")
}
