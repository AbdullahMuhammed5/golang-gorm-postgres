package controllers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	writer := makeRequest("GET", "/api/users/me", nil, &TestingConfigs.userToken)
	assert.Equal(t, http.StatusOK, writer.Code)
}
