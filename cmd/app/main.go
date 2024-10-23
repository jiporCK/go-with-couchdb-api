package main

import (
	"e-learning/go-with-couchdb/internal/controller"
	"e-learning/go-with-couchdb/internal/database"
	"e-learning/go-with-couchdb/internal/repository"
	"e-learning/go-with-couchdb/internal/usecase"
	"e-learning/go-with-couchdb/routes"
)

func main() {
	
	database.InitDB()

	productRepo := &repository.ProductRepo{}
	productService := usecase.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	router := routes.InitRoutes(productController)

	router.Run(":8081")
}

