package db

import (
	"CarServiceCenter/src/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// ClientInterface interface
type ClientInterface interface {
	OpenConnection() *mongo.Client
	CreateAppointment(models.Appointment) *models.Appointment
	DeleteAppointment(string) bool
	GetAppointment(string) (*models.Appointment, error)
	GetAppointmentsWithinDateRange(time.Time, time.Time) *[]models.Appointment
	UpdateAppointmentStatus(string, string) bool
}

// MongoStruct - implements ClientInterface
type MongoStruct struct{}

// OpenConnection - connects to local mongodb instance
func (d *MongoStruct) OpenConnection() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

// CreateAppointment - writes to db to store appointment and returns the created appointment
func (d *MongoStruct) CreateAppointment(appointment models.Appointment) *models.Appointment {
	client := d.OpenConnection()
	collection := client.Database("test").Collection("appointments")

	insertResult, err := collection.InsertOne(context.TODO(), appointment)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	appointment.ID = insertResult.InsertedID.(primitive.ObjectID)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return &appointment
}

// DeleteAppointment - deletes an appointment by given id and returns true if successful
func (d *MongoStruct) DeleteAppointment(appointmentID string) bool {
	response := true
	client := d.OpenConnection()
	collection := client.Database("test").Collection("appointments")

	objectID, err := primitive.ObjectIDFromHex(appointmentID)
	if err != nil {
		log.Println("DeleteAppointment: couldn't convert appointment ID from input:", err)
		response = false
	}
	documentID := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.TODO(), documentID)
	if err != nil {
		log.Println("DeleteAppointment: couldn't delete appointment from db:", err)
		response = false
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
	return response
}

// UpdateAppointmentStatus - writes to db to update specified appointment with new status
func (d *MongoStruct) UpdateAppointmentStatus(appointmentID, newStatus string) bool {
	client := d.OpenConnection()
	collection := client.Database("test").Collection("appointments")
	response := true

	objectID, err := primitive.ObjectIDFromHex(appointmentID)
	if err != nil {
		log.Println("UpdateAppointmentStatus: couldn't convert appointment ID from input:", err)
		response = false
	}
	var appointment models.Appointment
	documentID := bson.M{"_id": objectID}
	err = collection.FindOne(context.TODO(), documentID).Decode(&appointment)
	if err != nil {
		log.Println("UpdateAppointmentStatus: couldn't decode appointment from db:", err)
		response = false
	}
	appointment.Status = newStatus

	_, err = collection.UpdateOne(
		context.TODO(),
		documentID,
		bson.M{
			"$set": bson.M{"status": appointment.Status},
		},
	)
	if err != nil {
		log.Println("UpdateAppointmentStatus: unable to update status:", err)
		response = false
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return response
}

// GetAppointment - returns appointment by provided ID
func (d *MongoStruct) GetAppointment(appointmentID string) (*models.Appointment, error) {
	client := d.OpenConnection()
	collection := client.Database("test").Collection("appointments")
	objectID, err := primitive.ObjectIDFromHex(appointmentID)
	if err != nil {
		log.Println("GetAppointment: couldn't convert appointment ID from input:", err)
	}

	var result models.Appointment
	dbErr := collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if dbErr != nil {
		dbErr = fmt.Errorf("Unable to retrieve appointment with ID %v", appointmentID)
		log.Println("GetAppointment:", dbErr)
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &result, dbErr
}

// GetAppointmentsWithinDateRange - queries database for all appointments with dates that fall between given start and end dates and returns as list
func (d *MongoStruct) GetAppointmentsWithinDateRange(start, end time.Time) *[]models.Appointment {
	client := d.OpenConnection()
	collection := client.Database("test").Collection("appointments")

	var results []models.Appointment
	cur, err := collection.Find(context.TODO(), bson.M{"date": bson.M{"$gte": start, "$lte": end}})
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var appointment models.Appointment
		err := cur.Decode(&appointment)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, appointment)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &results
}
