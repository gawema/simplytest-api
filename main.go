package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"simplytest-api/api"
	"simplytest-api/storage/mongodb"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize MongoDB connection
	db, err := mongodb.NewMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	// Ensure MongoDB connection is closed when the program exits
	defer db.Close()

	// Create a new router using gorilla/mux
	r := mux.NewRouter()

	// Setup routes
	api.SetupRoutes(r, db.Collection())

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")
}
