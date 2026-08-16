package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/domain/entity"
	"github.com/go-inventory-v2/application/infrastructure/repository"
)

type ProductUsecase struct {
	productRepository repository.IProductRepository
}

type IProductUseCase interface {
	ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductInventoryPut(ctx context.Context, product entity.Product) (*entity.Product, error)
}

func NewProductUseCase(productRepository repository.IProductRepository) IProductUseCase {
	logger.InfoOutCtx("initializing product usecase SUCCESSFULLY")

	return &ProductUsecase{
		productRepository: productRepository,
	}
}

func (p *ProductUsecase) ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product usecase ProductAdd called")

	res, err := p.productRepository.ProductAdd(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product usecase ProductAdd failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductUsecase) ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product usecase ProductGet called")

	res, err := p.productRepository.ProductGet(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product usecase ProductGet failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductUsecase) ProductInventoryPut(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product usecase ProductInventoryPut called")

	res, err := p.productRepository.ProductInventoryPut(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product usecase ProductInventoryPut failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}