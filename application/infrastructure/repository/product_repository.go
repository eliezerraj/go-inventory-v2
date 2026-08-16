package repository

import (
	"context"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/domain/entity"
)

type ProductRepository struct {
}

type IProductRepository interface {
	ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductInventoryPut(ctx context.Context, product entity.Product) (*entity.Product, error)
}

func NewProductRepository() IProductRepository {
	logger.InfoOutCtx("initializing product repository SUCCESSFULLY")

	return &ProductRepository{}
}

func (p *ProductRepository) ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product repository ProductAdd called")

	return &product, nil
}

func (p *ProductRepository) ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product repository ProductGet called")

	return &product, nil
}

func (p *ProductRepository) ProductInventoryPut(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product repository ProductInventoryPut called")

	return &product, nil
}