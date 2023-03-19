package main

import (
	"log"
	"net/http"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/controllers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/initializers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	App initializers.App
)

func init() {
	// config, err := initializers.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("? Could not load environment variables", err)
	// }

	// initializers.ConnectDB(&config)

	App = initializers.InitializeEnv(".", "app")

	AuthController = controllers.NewAuthController(App.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(App.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	server = gin.Default()
}

func main() {
	// config, err := initializers.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("? Could not load environment variables", err)
	// }

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", App.Config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + App.Config.ServerPort))
}
