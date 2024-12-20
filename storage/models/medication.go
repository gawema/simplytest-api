package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Medication defines the data structure for a medication item
type Medication struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}
