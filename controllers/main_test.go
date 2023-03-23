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

	TicketControllerTest TicketController
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
	TicketControllerTest = NewTicketController(App.DB)

	router := gin.Default()
	routes := router.Group("/api")
	routes.POST("/auth/register", AuthControllerTest.SignUpUser)
	routes.POST("/auth/login", AuthControllerTest.SignInUser)
	routes.GET("/auth/refresh", AuthControllerTest.RefreshAccessToken)

	routes.GET("/users/me", middleware.DeserializeUser(), UserControllerTest.GetMe)

	routes.Use(middleware.DeserializeUser())
	routes.POST("/tickets", TicketControllerTest.CreateTicket)
	routes.GET("/tickets", TicketControllerTest.FindTickets)
	routes.PATCH("/tickets/:ticketId", TicketControllerTest.UpdateTicket)
	routes.GET("/tickets/:ticketId", TicketControllerTest.FindTicketById)
	routes.DELETE("/tickets/:ticketId", TicketControllerTest.DeleteTicket)
	routes.PATCH("/tickets/:ticketId/status", middleware.OnlyAdmin(), TicketControllerTest.UpdateTicketStatus)

	return router
}

func setup() {
	App = initializers.InitializeEnv("../", "app.test")
	App.DB.Exec(`
		CREATE TYPE ticket_status AS ENUM (
			'NEW',
			'IN_PROGRESS',
			'RESOLVED'
		);
	`)
	App.DB.AutoMigrate(&models.User{}, &models.Ticket{})
	newUser := models.SignUpInput{
		Name:            "Test",
		Email:           "test@email.com",
		Password:        "test1234",
		PasswordConfirm: "test1234",
	}
	makeRequest("POST", "/api/auth/register", newUser, false)

	// register an admin user
	newAdmin := models.SignUpInput{
		Name:            "admin",
		Email:           "admin@email.com",
		Password:        "test1234",
		PasswordConfirm: "test1234",
	}
	makeRequest("POST", "/api/auth/register", newAdmin, false)
	// make a user admin manually so we can use it in some tests
	initializers.AppInstance.DB.Model(&models.User{}).Where("email = ?", "admin@email.com").Update("role", "admin")
}

func teardown() {
	migrator := App.DB.Migrator()
	migrator.DropTable(&models.User{}, &models.Ticket{})
	App.DB.Exec(`DROP TYPE ticket_status;`)
}

func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+bearerToken("test@email.com"))
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func makeAdminRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+bearerToken("admin@email.com"))
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func bearerToken(email string) string {
	user := models.SignInInput{
		Email:    email,
		Password: "test1234",
	}

	writer := makeRequest("POST", "/api/auth/login", user, false)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["access_token"]
}

// func bearerToken() string {
// 	user := models.SignUpInput{
// 		Name:            "Test",
// 		Email:           "testuser@email.com",
// 		Password:        "test1234",
// 		PasswordConfirm: "test1234",
// 	}

// 	writer := makeRequest("POST", "/api/auth/register", user, false)
// 	var response map[string]interface{}
// 	json.Unmarshal(writer.Body.Bytes(), &response)
// 	data, _ := response["data"].(map[string]interface{})
// 	token, _ := data["access_token"].(string)
// 	return token
// }
