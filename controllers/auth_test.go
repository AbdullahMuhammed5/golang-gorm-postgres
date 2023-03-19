package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	newUser := models.SignUpInput{
		Name:            "abdullah",
		Email:           "abdullah@email.com",
		Password:        "test1234",
		PasswordConfirm: "test1234",
	}
	writer := makeRequest("POST", "/api/auth/register", newUser, false)
	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestLogin(t *testing.T) {
	user := models.SignInInput{
		Email:    "abdullah@email.com",
		Password: "test1234",
	}

	writer := makeRequest("POST", "/api/auth/login", user, false)

	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	_, exists := response["access_token"]

	assert.Equal(t, true, exists)
}

// func TestRefreshToken(t *testing.T) {
// 	writer := makeRequest("GET", "/api/auth/refresh", nil, true)

// 	assert.Equal(t, http.StatusOK, writer.Code)

// 	var response map[string]string
// 	json.Unmarshal(writer.Body.Bytes(), &response)
// 	_, exists := response["access_token"]

// 	assert.Equal(t, true, exists)
// }
