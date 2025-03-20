package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"e-learning/go-with-couchdb/internal/database"
	"e-learning/go-with-couchdb/internal/entity"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
)

type ProductRepo struct{}

// CreateProduct creates a new product, ensuring the name is unique
func (r *ProductRepo) CreateProduct(ctx context.Context, product entity.Product) error {
	db := database.GetDB("ishopdb")

	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	// Check if a product with the same name already exists
	exists, err := r.CheckProductNameExists(ctx, product.Name, "")
	if err != nil {
		log.Println("Error checking product name:", err)
		return fmt.Errorf("failed to check product name: %w", err)
	}
	if exists {
		log.Println("Product with name already exists:", product.Name)
		return fmt.Errorf("product with name '%s' already exists", product.Name)
	}

	_, err = db.Put(ctx, product.ID, product)
	if err != nil {
		log.Println("Database error:", err)
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetAllProducts retrieves all products from the database
func (r *ProductRepo) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	db := database.GetDB("ishopdb")

	rows, err := db.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		log.Println("Failed to retrieve products:", err)
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}

	var products []entity.Product
	for rows.Next() {
		// Get the document ID to check if it's a design document
		docID := rows.ID()
		if strings.HasPrefix(docID, "_design/") {
			continue // Skip design documents
		}

		var product entity.Product
		if err := rows.ScanDoc(&product); err != nil {
			log.Println("Failed to scan product:", err)
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

// GetProductById retrieves a product by its ID
func (r *ProductRepo) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	db := database.GetDB("ishopdb")

	row := db.Get(ctx, id)
	if err := row.Err; err != nil {
		if kivik.StatusCode(err) == 404 { // Not Found
			return nil, fmt.Errorf("product with ID %s not found", id)
		}
		log.Println("Failed to retrieve product:", err)
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	var product entity.Product
	if err := row.ScanDoc(&product); err != nil {
		log.Println("Failed to scan product document:", err)
		return nil, fmt.Errorf("failed to scan product document: %w", err)
	}

	return &product, nil
}

// UpdateProductById updates an existing product by ID
func (r *ProductRepo) UpdateProductById(ctx context.Context, id string, updatedProduct entity.Product) error {
	db := database.GetDB("ishopdb")

	// Fetch the existing product
	row := db.Get(ctx, id)
	if err := row.Err; err != nil {
		if kivik.StatusCode(err) == 404 {
			return fmt.Errorf("product with ID %s not found", id)
		}
		log.Println("Failed to retrieve product:", err)
		return fmt.Errorf("failed to retrieve product: %w", err)
	}

	var existingProduct entity.Product
	if err := row.ScanDoc(&existingProduct); err != nil {
		log.Println("Failed to scan product document:", err)
		return fmt.Errorf("failed to scan product document: %w", err)
	}

	// Check for revision mismatch
	if updatedProduct.Rev != existingProduct.Rev {
		log.Println("Document revision mismatch. Please try again")
		return fmt.Errorf("revision mismatch: expected %s, got %s", existingProduct.Rev, updatedProduct.Rev)
	}

	// Check if the updated name already exists (exclude current product)
	if updatedProduct.Name != existingProduct.Name {
		exists, err := r.CheckProductNameExists(ctx, updatedProduct.Name, id)
		if err != nil {
			log.Println("Error checking product name:", err)
			return fmt.Errorf("failed to check product name: %w", err)
		}
		if exists {
			log.Println("Product with name already exists:", updatedProduct.Name)
			return fmt.Errorf("product with name '%s' already exists", updatedProduct.Name)
		}
	}

	// Update fields
	existingProduct.Name = updatedProduct.Name
	existingProduct.Price = updatedProduct.Price

	// Save the updated product
	_, err := db.Put(ctx, id, existingProduct)
	if err != nil {
		log.Println("Failed to update product:", err)
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// DeleteProductById deletes a product by its ID and revision
func (r *ProductRepo) DeleteProductById(ctx context.Context, id string, rev string) error {
	db := database.GetDB("ishopdb")

	_, err := db.Delete(ctx, id, rev)
	if err != nil {
		if kivik.StatusCode(err) == 404 {
			return fmt.Errorf("product with ID %s not found", id)
		}
		log.Println("Failed to delete product:", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// BulkCreateProducts creates multiple products in a single operation
func (r *ProductRepo) BulkCreateProducts(ctx context.Context, products []entity.Product) error {
	db := database.GetDB("ishopdb")
	var docs []interface{}

	// Validate names for all products
	for _, product := range products {
		exists, err := r.CheckProductNameExists(ctx, product.Name, "")
		if err != nil {
			log.Println("Error checking product name:", err)
			return fmt.Errorf("failed to check product name: %w", err)
		}
		if exists {
			log.Println("Product with name already exists:", product.Name)
			return fmt.Errorf("product with name '%s' already exists", product.Name)
		}
	}

	// Prepare documents
	for i, product := range products {
		if product.ID == "" {
			products[i].ID = uuid.New().String()
		}
		docs = append(docs, products[i])
	}

	_, err := db.BulkDocs(ctx, docs)
	if err != nil {
		log.Println("Failed to create products in bulk:", err)
		return fmt.Errorf("failed to create products in bulk: %w", err)
	}

	return nil
}

// BulkUpdateProducts updates multiple products in a single operation
func (r *ProductRepo) BulkUpdateProducts(ctx context.Context, products []entity.Product) error {
	db := database.GetDB("ishopdb")
	var docs []interface{}

	// Validate names and prepare documents
	for _, product := range products {
		if product.ID == "" || product.Rev == "" {
			log.Println("Product ID or Rev missing, skipping update for product:", product)
			continue
		}

		// Fetch the existing product to get the current name
		existing, err := r.GetProductById(ctx, product.ID)
		if err != nil {
			log.Println("Failed to fetch product for validation:", err)
			continue
		}

		// Check if the name has changed and validate uniqueness
		if product.Name != existing.Name {
			exists, err := r.CheckProductNameExists(ctx, product.Name, product.ID)
			if err != nil {
				log.Println("Error checking product name:", err)
				continue
			}
			if exists {
				log.Println("Product with name already exists:", product.Name)
				continue
			}
		}

		docs = append(docs, product)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid products to update")
	}

	_, err := db.BulkDocs(ctx, docs)
	if err != nil {
		log.Println("Failed to update products in bulk:", err)
		return fmt.Errorf("failed to update products in bulk: %w", err)
	}

	return nil
}

// CheckProductNameExists checks if a product with the given name already exists
func (r *ProductRepo) CheckProductNameExists(ctx context.Context, name string, excludeID string) (bool, error) {
	db := database.GetDB("ishopdb")
	rows, err := db.Query(ctx, "_design/products", "_view/by_name", kivik.Options{
		"key": name, // Exact match for the name
	})
	if err != nil {
		log.Println("Failed to query products by name:", err)
		return false, fmt.Errorf("failed to query products by name: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.ScanValue(&id); err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}
		// If the found product's ID differs from excludeID, it’s a duplicate
		if id != excludeID {
			return true, nil
		}
	}
	return false, nil
}