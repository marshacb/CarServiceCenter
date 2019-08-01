package controller

import (
	"CarServiceCenter/src/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBTestImplementation struct{}

func (d *DBTestImplementation) OpenConnection() *mongo.Client {
	fmt.Println("being called")
	return nil
}

func (d *DBTestImplementation) CreateAppointment(appointment models.Appointment) *models.Appointment {
	date, _ := time.Parse(time.RFC3339, "2019-08-28T09:00:01+00:00")
	return &models.Appointment{
		Name:        "Test",
		Date:        date,
		Description: "Test Appointment",
		Status:      "open",
	}
}
func (d *DBTestImplementation) DeleteAppointment(id string) bool {
	if id != "1" {
		return false
	}
	return true
}
func (d *DBTestImplementation) GetAppointment(id string) (*models.Appointment, error) {
	date, _ := time.Parse(time.RFC3339, "2019-08-28T09:00:01+00:00")
	return &models.Appointment{
		Name:        "Test",
		Date:        date,
		Description: "Test Appointment",
		Status:      "open",
	}, nil
}
func (d *DBTestImplementation) GetAppointmentsWithinDateRange(start, end time.Time) *[]models.Appointment {
	fmt.Println("times", start, end)
	date, _ := time.Parse(time.RFC3339, "2019-08-28T09:00:01+00:00")
	return &[]models.Appointment{
		models.Appointment{
			Name:        "Test",
			Date:        date,
			Description: "Test Appointment",
			Status:      "open",
		},
		models.Appointment{
			Name:        "Test2",
			Date:        date,
			Description: "Test2 Appointment",
			Status:      "open",
		},
	}
}
func (d *DBTestImplementation) UpdateAppointmentStatus(id, status string) bool {
	if id == "2" {
		return false
	}
	return true
}

func TestCreateAppointmentSuccess(t *testing.T) {
	requestBody := map[string]interface{}{
		"Name":        "Ultimate Car Appointment",
		"Description": "even newer engine appointment",
		"Status":      "open",
		"Date":        "2019-08-28T09:00:01+00:00",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/appointment", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.CreateAppointment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"000000000000000000000000","name":"Test","description":"Test Appointment","status":"open","date":"2019-08-28T09:00:01Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestBadCreateAppointment(t *testing.T) {
	requestBody := map[string]interface{}{
		"Name":        "Ultimate Car Appointment",
		"Description": "",
		"Status":      "open",
		"Date":        "2019-08-28T09:00:01+00:00",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/appointment", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.CreateAppointment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "appointment must have valid name, description and date values"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteAppointmentSuccess(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/appointment/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.DeleteAppointment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "appointment 1 successfully deleted"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestBadDeleteAppointment(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/appointment/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.DeleteAppointment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "unable to find resource with id 2"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdateAppointmentStatusSuccess(t *testing.T) {
	requestBody := map[string]interface{}{
		"status": "closed",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("Patch", "/appointment/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.UpdateAppointmentStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "appointment status successfully updated to {closed}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestBadUpdateAppointmentStatus(t *testing.T) {
	requestBody := map[string]interface{}{
		"status": "closed",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("Patch", "/appointment/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appointmentsController.UpdateAppointmentStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "unable to update appointment status at id 2"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestBadGetAppointmentsWithinDateRange(t *testing.T) {
	req, err := http.NewRequest("GET", "/appointments/range", nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("start", "2019-07-29T09:00:01+00:00")
	rctx.URLParams.Add("end", "2019-08-29T09:00:01+00:00")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	appointmentsController := AppointmentsController{DB: &DBTestImplementation{}}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(appointmentsController.GetAppointmentsWithinDateRange)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "request must have valid start and end date range"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
