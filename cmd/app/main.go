package main

import (
	"e-learning/go-with-couchdb/internal/controller"
	"e-learning/go-with-couchdb/internal/database"
	"e-learning/go-with-couchdb/internal/repository"
	"e-learning/go-with-couchdb/internal/usecase"
	"e-learning/go-with-couchdb/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	
	// Load environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Inject dependencies for product module
	productRepo := &repository.ProductRepo{}
	productService := usecase.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// Initialize routes and pass the ProductController
	router := routes.InitRoutes(productController)

	// Start server on port 8081
	router.Run(":8081")
}
