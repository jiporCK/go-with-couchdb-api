package controller

import (
	"net/http"
	"fmt"
	"e-learning/go-with-couchdb/internal/entity"
	"e-learning/go-with-couchdb/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	service  *usecase.ProductService
	validate *validator.Validate
}

func NewProductController(s *usecase.ProductService) *ProductController {
	return &ProductController{
		service:  s,
		validate: validator.New(),
	}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := c.validate.Struct(product); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = fieldError.Error()
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errorMessages,
		})
		return
	}

	// Pass the request context to the service
	if err := c.service.CreateProduct(ctx.Request.Context(), product); err != nil {
		// Check for specific errors (e.g., duplicate name)
		if err.Error() == fmt.Sprintf("product with name '%s' already exists", product.Name) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	product, err := c.service.GetProductById(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %s not found", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product": product})
}

func (c *ProductController) UpdateProductById(ctx *gin.Context) {
	id := ctx.Param("_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	// Fetch the existing product to get the current revision
	existingProduct, err := c.service.GetProductById(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %s not found", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		return
	}

	var updatedProduct entity.Product
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate the updated product
	if err := c.validate.Struct(updatedProduct); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = fieldError.Error()
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errorMessages,
		})
		return
	}

	// Update the product
	err = c.service.UpdateProductById(ctx.Request.Context(), id, updatedProduct)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with name '%s' already exists", updatedProduct.Name) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == fmt.Sprintf("revision mismatch: expected %s, got %s", existingProduct.Rev, updatedProduct.Rev) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Revision mismatch, please refresh and try again"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		return
	}

	// Fetch the updated product to return the latest revision
	updated, err := c.service.GetProductById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": updated})
}

func (c *ProductController) DeleteProductById(ctx *gin.Context) {
	id := ctx.Param("_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	// Fetch the existing product to get the current revision
	existingProduct, err := c.service.GetProductById(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %s not found", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		return
	}

	// Delete the product
	err = c.service.DeleteProductById(ctx.Request.Context(), id, existingProduct.Rev)
	if err != nil {
		if err.Error() == fmt.Sprintf("product with ID %s not found", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (c *ProductController) BulkCreateProducts(ctx *gin.Context) {
	var products []entity.Product
	if err := ctx.ShouldBindJSON(&products); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate each product
	for i, product := range products {
		if err := c.validate.Struct(product); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMessages := make(map[string]string)
			for _, fieldError := range validationErrors {
				errorMessages[fieldError.Field()] = fieldError.Error()
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   fmt.Sprintf("Validation failed for product at index %d", i),
				"details": errorMessages,
			})
			return
		}
	}

	if err := c.service.BulkCreateProducts(ctx.Request.Context(), products); err != nil {
		if err.Error() == "no valid products to create" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid products to create"})
			return
		}
		// Check for duplicate name errors
		for _, product := range products {
			if err.Error() == fmt.Sprintf("product with name '%s' already exists", product.Name) {
				ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create products in bulk: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Products created successfully", "products": products})
}

func (c *ProductController) BulkUpdateProducts(ctx *gin.Context) {
	var products []entity.Product
	if err := ctx.ShouldBindJSON(&products); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate each product and ensure ID and Rev are present
	for i, product := range products {
		if product.ID == "" || product.Rev == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Product at index %d is missing ID or Rev", i),
			})
			return
		}
		if err := c.validate.Struct(product); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMessages := make(map[string]string)
			for _, fieldError := range validationErrors {
				errorMessages[fieldError.Field()] = fieldError.Error()
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   fmt.Sprintf("Validation failed for product at index %d", i),
				"details": errorMessages,
			})
			return
		}
	}

	if err := c.service.BulkUpdateProducts(ctx.Request.Context(), products); err != nil {
		if err.Error() == "no valid products to update" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid products to update"})
			return
		}
		// Check for duplicate name errors
		for _, product := range products {
			if err.Error() == fmt.Sprintf("product with name '%s' already exists", product.Name) {
				ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update products in bulk: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Products updated successfully"})
}