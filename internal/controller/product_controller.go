package controller

import (
	"e-learning/go-with-couchdb/internal/entity"
	"e-learning/go-with-couchdb/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	service *usecase.ProductService
	validate *validator.Validate
}

func NewProductController(s *usecase.ProductService) *ProductController {
	return &ProductController{
		service: s,
		validate: validator.New(),
	}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreateProduct(product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("_id")

	product, err := c.service.GetProductById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error: ": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) UpdateProductById(ctx *gin.Context) {

	id := ctx.Param("_id")
	product, err := c.service.GetProductById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error: ": "Product not found"})
		return
	}

	var updatedProduct entity.Product
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error: ": "Invalid input"})
		return
	}
	updateProduct := c.service.UpdateProductById(id, product.Rev, updatedProduct)
	if updateProduct != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, updateProduct)
}

func (c *ProductController) DeleteProductById(ctx *gin.Context) {

	id := ctx.Param("_id")
	existingProduct, err := c.service.GetProductById(id)
    if err != nil {
        ctx.JSON(404, gin.H{"error": "Product not found"})
        return
    }
	c.service.DeleteProductById(id, existingProduct.Rev)
	ctx.JSON(204, gin.H{"message": "Product deleted successfully"})
}

func (c *ProductController) BulkCreateProducts(ctx *gin.Context) {
	var products []entity.Product

	if err := ctx.ShouldBindJSON(&products); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.BulkCreateProducts(products); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create products in bulk"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Products created successfully"})
}

func (c *ProductController) BulkUpdateProducts(ctx *gin.Context) {
	var products []entity.Product

	if err := ctx.ShouldBindJSON(&products); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, product := range products {
		if product.ID == "" || product.Rev == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID and Rev are required for each product"})
			return
		}
	}

	if err := c.service.BulkUpdateProducts(products); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update products in bulk"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Products updated successfully"})
}