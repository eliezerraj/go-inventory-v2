package controller

import (
	"context"

	"github.com/eliezerraj/go-core/v3/logger"
	"go.uber.org/zap"
	"github.com/go-inventory-v2/application/domain/usecase"
	"github.com/go-inventory-v2/application/domain/external"
	"github.com/go-inventory-v2/application/domain/entity"
)

type ProductController struct {
	productUseCase usecase.IProductUseCase
}

func NewProductController(productUseCase usecase.IProductUseCase) *ProductController {
	logger.InfoOutCtx("initializing product controller SUCCESSFULLY")

	return &ProductController{
		productUseCase: productUseCase,
	}
}

func (p *ProductController) ProductAdd(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.InfoOutCtx("product controller ProductAdd called")

	product := entity.Product{
		Name:        req.Name,
	}

	res, err := p.productUseCase.ProductAdd(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductAdd failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductGet(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.InfoOutCtx("product controller ProductGet called")

	product := entity.Product{
		Name:        req.Name,
	}

	res, err := p.productUseCase.ProductGet(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductGet failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductInventoryPut(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.InfoOutCtx("product controller ProductInventoryPut called")

	product := entity.Product{
		Name:        req.Name,
	}

	res, err := p.productUseCase.ProductInventoryPut(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductInventoryPut failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}