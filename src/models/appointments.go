package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Appointment - type that represents a users appointment
type Appointment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	Date        time.Time          `json:"date" bson:"date"`
}

// Status - status for appointment in update
type Status struct {
	Status string `json:"status" bson:"status"`
}

// AppointmentUpdate - type representing appointment update
type AppointmentUpdate struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status string             `json:"status" bson:"status"`
}
