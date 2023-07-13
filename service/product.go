package service

import (
	"context"
	"go-jwt/domain"
)

type productService struct {
	productRepo domain.ProductRepository
}

func NewProductService(p domain.ProductRepository) domain.ProductService {
	return productService{productRepo: p}
}

func (p productService) GetProductByName(ctx context.Context, name string) (domain.Products, error) {
	return p.productRepo.GetByName(ctx, name)
}
