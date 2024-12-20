package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Medication struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var medications = []Medication{
	{ID: 1, Name: "Paracetamol", Description: "Pain reliever and fever reducer."},
	{ID: 2, Name: "Ibuprofen", Description: "Reduces inflammation and treats pain."},
}

var nextID = 3

func getMedications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medications)
}

func getMedicationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	for _, medication := range medications {
		if medication.ID == id {
			json.NewEncoder(w).Encode(medication)
			return
		}
	}

	http.Error(w, "Medication not found", http.StatusNotFound)
}

func createMedication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMedication Medication
	if err := json.NewDecoder(r.Body).Decode(&newMedication); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newMedication.ID = nextID
	nextID++
	medications = append(medications, newMedication)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMedication)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/medications", getMedications).Methods("GET")
	r.HandleFunc("/medications/{id}", getMedicationByID).Methods("GET")
	r.HandleFunc("/medications", createMedication).Methods("POST")

	http.ListenAndServe(":8080", r)
}
