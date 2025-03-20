package main

import (
	"e-learning/go-with-couchdb/internal/controller"
	"e-learning/go-with-couchdb/internal/database"
	"e-learning/go-with-couchdb/internal/repository"
	"e-learning/go-with-couchdb/internal/usecase"
	"e-learning/go-with-couchdb/routes"
)

func main() {
	
	// Initialize CouchDB database connection
	database.InitDB()

	// Inject dependencies for product module
	productRepo := &repository.ProductRepo{}
	productService := usecase.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// Initialize routes and pass the ProductController
	router := routes.InitRoutes(productController)

	// Start server on port 8081
	router.Run(":8081")
}
