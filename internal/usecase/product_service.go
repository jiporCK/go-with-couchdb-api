package usecase

import (
	"e-learning/go-with-couchdb/internal/entity"
	"e-learning/go-with-couchdb/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product entity.Product) error {
    return s.repo.CreateProduct(product)
}

func (s *ProductService) GetAllProducts() ([]entity.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) GetProductById(id string) (*entity.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *ProductService) UpdateProductById(id string, rev string, updatedProduct entity.Product) error {
	return s.repo.UpdateProductById(id, rev, updatedProduct)
}

func (s *ProductService) DeleteProductById(id string, rev string) error {
	return s.repo.DeleteProductById(id, rev)
}

func (s *ProductService) BulkCreateProducts(products []entity.Product) error {
	return s.repo.BulkCreateProducts(products)
}

func (s *ProductService) BulkUpdateProducts(products []entity.Product) error {
	return s.repo.BulkUpdateProducts(products)
}
