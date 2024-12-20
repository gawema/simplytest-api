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

	// Define routes for medication CRUD operations
	// Each route maps to a specific HTTP method and handler function
	r.HandleFunc("/medications", medicationHandler.GetMedications).Methods("GET")           // List all medications
	r.HandleFunc("/medications/{id}", medicationHandler.GetMedicationByID).Methods("GET")   // Get one medication
	r.HandleFunc("/medications", medicationHandler.CreateMedication).Methods("POST")        // Create medication
	r.HandleFunc("/medications/{id}", medicationHandler.UpdateMedication).Methods("PUT")    // Update medication
	r.HandleFunc("/medications/{id}", medicationHandler.DeleteMedication).Methods("DELETE") // Delete medication
}
