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
	"github.com/abdullahmuhammed5/golang-gorm-postgres/utils"

	"github.com/gin-gonic/gin"
)

type testingConfigs struct {
	adminToken string
	userToken  string
}

var (
	server *gin.Engine
	App    initializers.App

	AuthControllerTest AuthController

	UserControllerTest UserController

	TicketControllerTest TicketController

	TestingConfigs testingConfigs
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

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

	hashedPassword, _ := utils.HashPassword("test1234")
	user := models.User{Name: "Test", Email: "test@email.com", Password: hashedPassword, Role: "user"}
	admin := models.User{Name: "Admin", Email: "admin@email.com", Password: hashedPassword, Role: "admin"}

	App.DB.Create(&user)
	App.DB.Create(&admin)

	TestingConfigs.adminToken = LoginAs("admin@email.com")
	TestingConfigs.userToken = LoginAs("test@email.com")
}

func teardown() {
	migrator := App.DB.Migrator()
	migrator.DropTable(&models.User{}, &models.Ticket{})
	App.DB.Exec(`DROP TYPE ticket_status;`)
}

func makeRequest(method, url string, body interface{}, accessToken *string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if accessToken != nil {
		request.Header.Add("Authorization", "Bearer "+*accessToken)
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

	writer := makeRequest("POST", "/api/auth/login", user, nil)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["access_token"]
}

func LoginAs(email string) string {
	user := models.SignInInput{
		Email:    email,
		Password: "test1234",
	}

	writer := makeRequest("POST", "/api/auth/login", user, nil)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["access_token"]
}
