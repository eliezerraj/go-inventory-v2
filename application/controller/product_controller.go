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
		Sku:         req.Sku,
		Name:        req.Name,
		Status:      req.Status,
		Type:        req.Type,
		LeadTime:    req.LeadTime,
	}

	if req.Price != nil {
		price := entity.Price{
			Currency:    req.Price.Currency,
			Amount:      req.Price.Amount,
		}
		product.Price = &price
	}

	if req.Inventory != nil {
		inventory := entity.Inventory{
			Available:   req.Inventory.Available,
			Pending:     req.Inventory.Pending,
			Sold:        req.Inventory.Sold,
		}
		product.Inventory = &inventory
	}
	
	res, err := p.productUseCase.ProductAdd(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductAdd failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductGet(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.Info(ctx, "product controller ProductGet called", zap.String("sku", req.Sku))

	product := entity.Product{
		Sku: req.Sku,
	}

	res, err := p.productUseCase.ProductGet(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductGet failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductPut(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.InfoOutCtx("product controller ProductPut called")

	product := entity.Product{
		Sku:         req.Sku,
		Name:        req.Name,
		Status:      req.Status,
		Type:        req.Type,
		LeadTime:    req.LeadTime,
	}

	if req.Price != nil {
		price := entity.Price{
			Currency:    req.Price.Currency,
			Amount:      req.Price.Amount,
		}
		product.Price = &price
	}

	if req.Inventory != nil {
		inventory := entity.Inventory{
			Available:   req.Inventory.Available,
			Pending:     req.Inventory.Pending,
			Sold:        req.Inventory.Sold,
		}
		product.Inventory = &inventory
	}

	res, err := p.productUseCase.ProductPut(ctx, product)
	if err != nil {
		logger.ErrorOutCtx("product controller ProductPut failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}