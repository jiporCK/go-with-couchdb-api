package routes

import (
	"e-learning/go-with-couchdb/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(controller *controller.ProductController) *gin.Engine {

	r := gin.Default()

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