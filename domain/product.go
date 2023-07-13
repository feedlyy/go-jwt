package domain

import (
	"context"
	"database/sql"
)

type Products struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Qty         int          `json:"qty"`
	Description string       `json:"description"`
	CreatedAt   sql.NullTime `json:"createdAt"`
	UpdatedAt   sql.NullTime `json:"updatedAt"`
}

type ProductRepository interface {
	GetByName(ctx context.Context, name string) (Products, error)
}

type ProductService interface {
	GetProductByName(ctx context.Context, name string) (Products, error)
}
