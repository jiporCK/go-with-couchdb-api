package routes

import (
	"e-learning/go-with-couchdb/internal/controller"

	"github.com/gin-gonic/gin"

	"log"
)

func InitRoutes(controller *controller.ProductController) *gin.Engine {

	// Create a new Gin router instance with default middleware
	r := gin.Default()

	// Set trusted proxies to only allow requests from specified IPs
	err := r.SetTrustedProxies([]string{"127.0.0.1", "192.168.0.0/16", "::1"})
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	// Create a group of routes related to products,
	productRouter := r.Group("/api/v1/products")
	{
		productRouter.POST("", controller.CreateProduct)
		productRouter.GET("", controller.GetAllProducts)
		productRouter.GET("/:_id", controller.GetProductById)
		productRouter.PUT("/:_id", controller.UpdateProductById)
		productRouter.DELETE("/:_id", controller.DeleteProductById)

		// For bulk create and update
		productRouter.POST("/bulk-create", controller.BulkCreateProducts)
		productRouter.PUT("/bulk-update", controller.BulkUpdateProducts)
	}

	return r
}