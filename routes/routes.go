package routes

import (
	"e-learning/go-with-couchdb/internal/controller"

	"github.com/gin-gonic/gin"

	"log"
)

func InitRoutes(controller *controller.ProductController) *gin.Engine {

	r := gin.Default()

	err := r.SetTrustedProxies([]string{"127.0.0.1", "192.168.0.0/16", "::1"})
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	productRouter := r.Group("/products")
	{
		productRouter.POST("", controller.CreateProduct)
		productRouter.GET("", controller.GetAllProducts)
		productRouter.GET("/:_id", controller.GetProductById)
		productRouter.PUT("/:_id", controller.UpdateProductById)
		productRouter.DELETE("/:_id", controller.DeleteProductById)

		productRouter.POST("/bulk-create", controller.BulkCreateProducts)
		productRouter.PUT("/bulk-update", controller.BulkUpdateProducts)
	}

	return r
}
