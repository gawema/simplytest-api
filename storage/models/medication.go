package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Medication represents a medication record in the database
// It uses struct tags to define how the struct should be serialized to JSON and BSON (MongoDB format)
type Medication struct {
	// ID is the unique identifier for the medication
	// omitempty means the field will be omitted from JSON/BSON if it's empty
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// Name of the medication
	Name string `json:"name" bson:"name"`

	// Description provides details about the medication
	Description string `json:"description" bson:"description"`

	// Price of the medication
	Price float64 `json:"price" bson:"price"`

	// ImageURL provides the URL to the medication's image
	ImageURL string `json:"imageUrl" bson:"imageUrl"`
}
