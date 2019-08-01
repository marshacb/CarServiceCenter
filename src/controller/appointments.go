package controller

import (
	"CarServiceCenter/src/db"
	"CarServiceCenter/src/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

// AppointmentsController - struct that has reference to db client
type AppointmentsController struct {
	DB db.ClientInterface
}

// CreateAppointment - accepts appointment name, description, and returns created appointment
func (a *AppointmentsController) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		log.Println("error decoding json", err.Error())
	}
	status := http.StatusOK
	response := []byte{}

	if len(appointment.Name) == 0 || appointment.Date.IsZero() || len(appointment.Description) == 0 {
		status = http.StatusBadRequest
		response = []byte(fmt.Sprintf("appointment must have valid name, description and date values"))
	} else {
		appointment.Status = "open"
		newAppointment := a.DB.CreateAppointment(appointment)
		response, err = json.Marshal(newAppointment)
		if err != nil {
			log.Println("error:", err)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// DeleteAppointment - accepts appointmentID to be deleted
func (a *AppointmentsController) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	status := http.StatusOK
	response := fmt.Sprintf("appointment %v successfully deleted", id)

	deleted := a.DB.DeleteAppointment(id)
	if deleted == false {
		status = http.StatusBadRequest
		response = fmt.Sprintf("unable to find resource with id %v", id)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// UpdateAppointmentStatus - accepts id and status to update appointment status
func (a *AppointmentsController) UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updatedStatus models.Status
	json.NewDecoder(r.Body).Decode(&updatedStatus)
	status := http.StatusOK
	response := fmt.Sprintf("appointment status successfully updated to %v", updatedStatus)

	updated := a.DB.UpdateAppointmentStatus(id, updatedStatus.Status)
	if updated == false {
		status = http.StatusBadRequest
		response = fmt.Sprintf("unable to update appointment status at id %v", id)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// GetAppointment - accepts appointment id and returns specified appointment
func (a *AppointmentsController) GetAppointment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	status := http.StatusOK
	response := []byte{}

	appointment, err := a.DB.GetAppointment(id)
	if err != nil {
		status = http.StatusBadRequest
		response = []byte(fmt.Sprintf("Unable to retrieve appointment with ID %v", id))
	} else {
		response, err = json.Marshal(appointment)
		if err != nil {
			log.Println("error marshaling appointment struct")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// GetAppointmentsWithinDateRange - accepts start and end date and returns all appointments within that range
func (a *AppointmentsController) GetAppointmentsWithinDateRange(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	response := []byte{}
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("error parsing url string", err)
	}
	urlValues := parsedURL.Query()
	var start, end time.Time
	if _, ok := urlValues["start"]; ok {
		startTime := strings.Replace(strings.Join(urlValues["start"], ""), " ", "+", -1)
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			log.Println("error parsing time for start")
		}
	}
	if _, ok := urlValues["end"]; ok {
		endTime := strings.Replace(strings.Join(urlValues["end"], ""), " ", "+", -1)
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			log.Println("err parsing time for end")
		}
	}

	if start.IsZero() || end.IsZero() || end.Before(start) {
		status = http.StatusBadRequest
		response = []byte(fmt.Sprintf("request must have valid start and end date range"))
	} else {
		results := a.DB.GetAppointmentsWithinDateRange(start, end)
		response, err = json.Marshal(results)
		if err != nil {
			log.Println("error marshaling results")
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
