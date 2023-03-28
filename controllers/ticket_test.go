package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/initializers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_ticket_CreateTicket(t *testing.T) {
	ticket := models.CreateTicketRequest{
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	writer := makeRequest("POST", "/api/tickets", ticket, &TestingConfigs.userToken)

	assert.Equal(t, http.StatusCreated, writer.Code)

	var response map[string]interface{}
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	data, _ := response["data"].(map[string]interface{})
	if err != nil {
		t.Fatal(writer)
	}
	if data["id"] == 0 {
		t.Error("Expected ticket ID to be non-zero")
	}
}

func Test_ticket_CreateTicket_validationsFailed(t *testing.T) {
	ticket := models.CreateTicketRequest{
		Description: "Test Description",
	}

	writer := makeRequest("POST", "/api/tickets", ticket, &TestingConfigs.userToken)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func Test_ticket_UpdateTicket(t *testing.T) {
	// Seed test data
	ticket := &models.Ticket{
		Title:       "Ticket 1",
		Description: "Test Description",
		CreatedBy:   1,
		Status:      "NEW",
	}
	initializers.AppInstance.DB.Create(&ticket)

	updatedTicket := models.UpdateTicket{
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	writer := makeRequest("PATCH", "/api/tickets/"+strconv.FormatUint(uint64(ticket.ID), 10), updatedTicket, &TestingConfigs.userToken)

	assert.Equal(t, http.StatusOK, writer.Code)

	// Decode response body
	var response models.Response
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	data, _ := response.Data.(map[string]interface{})

	if err != nil {
		t.Fatal(writer)
	}
	if data["title"] != "Updated Title" || data["description"] != "Updated Description" {
		t.Error("Something went wrong!")
	}
}

func Test_ticket_OnlyAdminCanUpdateTicketStatus(t *testing.T) {
	// Seed test data
	ticket := &models.Ticket{
		Title:       "Ticket 1",
		Description: "Test Description",
		CreatedBy:   1,
		Status:      "NEW",
	}
	initializers.AppInstance.DB.Create(&ticket)

	updatedTicket := models.UpdateTicket{
		Status: "IN_PROGRESS",
	}

	url := "/api/tickets/" + strconv.FormatUint(uint64(ticket.ID), 10) + "/status"
	writer := makeRequest("PATCH", url, updatedTicket, &TestingConfigs.userToken)

	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func Test_ticket_UpdateTicketStatusToInProgress(t *testing.T) {
	// Seed test data
	ticket := &models.Ticket{
		Title:       "Ticket 1",
		Description: "Test Description",
		CreatedBy:   1,
		Status:      "NEW",
	}
	initializers.AppInstance.DB.Create(&ticket)

	updatedTicket := models.UpdateTicket{
		Status: "IN_PROGRESS",
	}

	url := "/api/tickets/" + strconv.FormatUint(uint64(ticket.ID), 10) + "/status"
	writer := makeRequest("PATCH", url, updatedTicket, &TestingConfigs.adminToken)

	assert.Equal(t, http.StatusOK, writer.Code)

	// Decode response body
	var response models.Response
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	data, _ := response.Data.(map[string]interface{})

	if err != nil {
		t.Fatal(writer)
	}
	if data["status"] != "IN_PROGRESS" {
		t.Error("Something went wrong!")
	}
}

func Test_ticket_UpdateTicketStatusToResolved(t *testing.T) {
	// Seed test data
	ticket := &models.Ticket{
		Title:       "Ticket 1",
		Description: "Test Description",
		CreatedBy:   1,
		Status:      "NEW",
	}
	initializers.AppInstance.DB.Create(&ticket)

	updatedTicket := models.UpdateTicket{
		Status: "RESOLVED",
	}

	url := "/api/tickets/" + strconv.FormatUint(uint64(ticket.ID), 10) + "/status"
	writer := makeRequest("PATCH", url, updatedTicket, &TestingConfigs.adminToken)

	assert.Equal(t, http.StatusOK, writer.Code)

	// Decode response body
	var response models.Response
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	data, _ := response.Data.(map[string]interface{})

	if err != nil {
		t.Fatal(writer)
	}
	if data["status"] != "RESOLVED" {
		t.Error("Something went wrong!")
	}
}

func Test_ticket_deleteTicket(t *testing.T) {
	// Seed test data
	ticket := &models.Ticket{
		Title:       "Ticket 1",
		Description: "Test Description",
		CreatedBy:   1,
		Status:      "NEW",
	}
	initializers.AppInstance.DB.Create(&ticket)

	writer := makeRequest("DELETE", "/api/tickets/"+strconv.FormatUint(uint64(ticket.ID), 10), nil, &TestingConfigs.userToken)

	assert.Equal(t, http.StatusOK, writer.Code)

	// Check that the ticket was deleted from the database
	var deletedTicket models.Ticket
	err := initializers.AppInstance.DB.First(&deletedTicket, ticket.ID).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err, "expected ticket to be deleted from the database")
}

func Test_ticket_ListTickets(t *testing.T) {
	// Seed test data
	initializers.AppInstance.DB.Exec("DELETE FROM tickets;")
	initializers.AppInstance.DB.Create(
		&models.Ticket{
			Title:       "Ticket 1",
			Description: "Test Description",
			CreatedBy:   1,
			Status:      "NEW",
		})
	initializers.AppInstance.DB.Create(
		&models.Ticket{
			Title:       "Ticket 2",
			Description: "Test Description",
			CreatedBy:   1,
			Status:      "NEW",
		})

	// Send GET request to /tickets endpoint
	writer := makeRequest("GET", "/api/tickets", nil, &TestingConfigs.userToken)

	// Check response status code
	assert.Equal(t, http.StatusOK, writer.Code)

	// Decode response body
	var tickets models.Response
	err := json.Unmarshal(writer.Body.Bytes(), &tickets)
	if err != nil {
		t.Fatal(err)
	}
	// Access the data field from the Response struct
	data, ok := tickets.Data.([]interface{})
	if !ok {
		t.Error("Data field is not an array")
	}

	// Check number of tickets
	assert.Equal(t, 2, len(data))
}
