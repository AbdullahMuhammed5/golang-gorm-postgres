package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/controllers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/routes"

	initializers "github.com/abdullahmuhammed5/golang-gorm-postgres/test/config"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func router() *gin.Engine {

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	server := gin.Default()

	router := server.Group("/api")

	AuthRouteController.AuthRoute(router)

	// publicRoutes := router.Group("/api/auth")
	// publicRoutes.POST("/register", AuthController.SignUpUser)
	// publicRoutes.POST("/login", AuthController.SignInUser)

	return server
}

func setup() {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	initializers.DB.AutoMigrate(&models.User{})
	// database.Database.AutoMigrate(&models.Entry{})
}

func teardown() {
	migrator := initializers.DB.Migrator()
	migrator.DropTable(&models.User{})
	// migrator.DropTable(&models.Entry{})
}

func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+bearerToken())
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func bearerToken() string {
	user := models.SignInInput{
		Email:    "abdullah@email.com",
		Password: "test1234",
	}

	writer := makeRequest("POST", "/api/auth/login", user, false)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["access_token"]
}
