package api

import (
	"simplytest-api/api/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(r *mux.Router, collection *mongo.Collection) {
	medicationHandler := handlers.NewMedicationHandler(collection)

	r.HandleFunc("/medications", medicationHandler.GetMedications).Methods("GET")
	r.HandleFunc("/medications/{id}", medicationHandler.GetMedicationByID).Methods("GET")
	r.HandleFunc("/medications", medicationHandler.CreateMedication).Methods("POST")
}
