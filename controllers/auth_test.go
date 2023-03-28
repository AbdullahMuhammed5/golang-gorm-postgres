package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/stretchr/testify/assert"
)

func Test_auth_register(t *testing.T) {
	newUser := models.SignUpInput{
		Name:            "abdullah",
		Email:           "abdullah@email.com",
		Password:        "test1234",
		PasswordConfirm: "test1234",
	}
	writer := makeRequestV1("POST", "/api/auth/register", newUser, nil)
	assert.Equal(t, http.StatusCreated, writer.Code)
}

func Test_auth_login(t *testing.T) {
	user := models.SignInInput{
		Email:    "abdullah@email.com",
		Password: "test1234",
	}

	writer := makeRequestV1("POST", "/api/auth/login", user, nil)

	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	_, exists := response["access_token"]

	assert.Equal(t, true, exists)
}

// func Test_auth_refresh_token(t *testing.T) {
// 	writer := makeRequestV1("GET", "/api/auth/refresh", nil, &TestingConfigs.userToken)

// 	assert.Equal(t, http.StatusOK, writer.Code)

// 	var response map[string]string
// 	json.Unmarshal(writer.Body.Bytes(), &response)
// 	_, exists := response["access_token"]

// 	assert.Equal(t, true, exists)
// }
