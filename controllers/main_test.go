package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/initializers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/middleware"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
	App    initializers.App

	AuthControllerTest AuthController
	// AuthRouteController routes.AuthRouteController

	UserControllerTest UserController
	// UserRouteController routes.UserRouteController
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

// func router() *gin.Engine {

// 	AuthController = controllers.NewAuthController(initializers.DB)
// 	AuthRouteController = routes.NewAuthRouteController(AuthController)

// 	server := gin.Default()

// 	router := server.Group("/api")

// 	AuthRouteController.AuthRoute(router)

// 	return server
// }

func router() *gin.Engine {
	AuthControllerTest = NewAuthController(App.DB)

	UserControllerTest = NewUserController(App.DB)

	router := gin.Default()
	publicRoutes := router.Group("/api")
	publicRoutes.POST("/auth/register", AuthControllerTest.SignUpUser)
	publicRoutes.POST("/auth/login", AuthControllerTest.SignInUser)
	publicRoutes.GET("/auth/refresh", AuthControllerTest.RefreshAccessToken)

	publicRoutes.GET("/users/me", middleware.DeserializeUser(), UserControllerTest.GetMe)

	return router
}

func setup() {
	App = initializers.InitializeEnv("../", "app.test")
	App.DB.AutoMigrate(&models.User{})
	// database.Database.AutoMigrate(&models.Entry{})
}

func teardown() {
	migrator := App.DB.Migrator()
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
