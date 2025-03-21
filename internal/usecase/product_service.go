package usecase

import (
	"context"
	"e-learning/go-with-couchdb/internal/entity"
	"e-learning/go-with-couchdb/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product entity.Product) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *ProductService) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	return s.repo.GetProductById(ctx, id)
}

func (s *ProductService) UpdateProductById(ctx context.Context, id string, product entity.Product) error {
	return s.repo.UpdateProductById(ctx, id, product)
}

func (s *ProductService) DeleteProductById(ctx context.Context, id string, rev string) error {
	return s.repo.DeleteProductById(ctx, id, rev)
}

func (s *ProductService) BulkCreateProducts(ctx context.Context, products []entity.Product) error {
	return s.repo.BulkCreateProducts(ctx, products)
}

func (s *ProductService) BulkUpdateProducts(ctx context.Context, products []entity.Product) error {
	return s.repo.BulkUpdateProducts(ctx, products)
}