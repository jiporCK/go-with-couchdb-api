package repository

import (
	"context"
	"e-learning/go-with-couchdb/internal/database"
	"e-learning/go-with-couchdb/internal/entity"
	"log"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"fmt"
)

type ProductRepo struct{}

func (r *ProductRepo) CreateProduct(product entity.Product) error {
    db := database.GetDB("products")
	
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

    _, err := db.Put(context.TODO(), product.ID, product)
    if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}	
    return nil
}


func (c *ProductRepo) GetAllProducts() ([]entity.Product, error) {
	db := database.GetDB("products")

	rows, err := db.AllDocs(context.TODO(), kivik.Options{"include_docs": true})

	if err != nil {
		log.Println("Failed to retrieve products", err)
		return nil, err
	}

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.ScanDoc(&product); err != nil {
			log.Println("Failed to scan products: ",err)
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

func (c *ProductRepo) GetProductById(id string) (*entity.Product, error) {
	db := database.GetDB("products")

	row := db.Get(context.TODO(), id)

	if row.Err != nil {
		log.Println("Failed to retrieve product: ", row.Err)
		return nil, row.Err
	}

	var product entity.Product
	if err := row.ScanDoc(&product); err != nil {
		log.Println("Failed to scan product document: ", err)
		return nil, err
	}

	return &product, nil
}

func (c *ProductRepo) UpdateProductById(id string, rev string, updatedProduct entity.Product) error {
	db := database.GetDB("products")

	row := db.Get(context.TODO(), id)
	if row.Err != nil {
        log.Println("Failed to retrieve product: ", row.Err)
        return nil
    }

	var existingProduct entity.Product

	if err := row.ScanDoc(&existingProduct); err != nil {
		log.Println("Failed to scan product document: ", err)
		return nil
	}

	if updatedProduct.Rev != existingProduct.Rev {
		log.Println("Document revision mismatch. Please try again")
		return fmt.Errorf("revision mismatch: expected %s, got %s", existingProduct.Rev, updatedProduct.Rev)
	}

	existingProduct.Name = updatedProduct.Name
	existingProduct.Price = updatedProduct.Price


	_, err := db.Put(context.TODO(), id, existingProduct)
	if err != nil {
		log.Println("Failed to update product: ", err)
		return nil
	}
	return nil
}

func (c *ProductRepo) DeleteProductById(id string, rev string) error {
	db := database.GetDB("products")

	_, err := db.Delete(context.TODO(), id, rev)
	if err != nil {
		log.Println("Failed to delete product: ", err)
		return err
	}
	return nil
}

func (c *ProductRepo) BulkCreateProducts(products []entity.Product) error {
	db := database.GetDB("products")
	var docs []interface{}

	for i, product := range products {
		if product.ID == "" {
			products[i].ID = uuid.New().String()
		}
		docs = append(docs, products[i])
	}

	_, err := db.BulkDocs(context.TODO(), docs)
	if err != nil {
		log.Println("Failed to create products in bulk: ", err)
		return err
	}

	return nil
}

func (c *ProductRepo) BulkUpdateProducts(products []entity.Product) error {
	db := database.GetDB("products")
	var docs []interface{}

	for _, product := range products {
		if product.ID == "" || product.Rev == "" {
			log.Println("Product ID or Rev missing, skipping update for product:", product)
			continue
		}
		docs = append(docs, product)
	}

	_, err := db.BulkDocs(context.TODO(), docs)
	if err != nil {
		log.Println("Failed to update products in bulk: ", err)
		return err
	}

	return nil
}

