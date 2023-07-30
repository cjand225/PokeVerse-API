package main

import (
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"pokeverse/web-service/controller"
	"pokeverse/web-service/database"
	"pokeverse/web-service/service"
)

func main() {
	// Create the Gin router with default middleware
	router := gin.Default()

	// Initialize the database connection pool
	pool := database.ConnectDatabase()

	// Create a new instance of the service with the connection pool
	service := service.NewService(pool)

	// Create a new instance of the controller with the service
	controller := controller.NewController(service)

	// Setup the routes using the controller
	controller.SetupRoutes(router)

	// Start the server and listen on port 8080
	err := router.Run("localhost:8080")

    if err != nil {
		log.Fatalf(err.Error())
        return
    }

	// Exit the application with status code 0 (success)
	os.Exit(0)
}
