package controller

import (
	"context"

	"go.opentelemetry.io/otel"

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
	tracer := otel.Tracer("product.controller")
	ctx, span := tracer.Start(ctx, "ProductController.ProductAdd")
	defer span.End()

	logger.Info(ctx, "product controller ProductAdd called")

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
		logger.Error(ctx, "product controller ProductAdd failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductGet(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	tracer := otel.Tracer("product.controller")
	ctx, span := tracer.Start(ctx, "ProductController.ProductGet")
	defer span.End()

	logger.Info(ctx, "product controller ProductGet called", zap.String("sku", req.Sku))

	product := entity.Product{
		ID:  req.ID,
		Sku: req.Sku,
	}

	res, err := p.productUseCase.ProductGet(ctx, product)
	if err != nil {
		logger.Error(ctx, "product controller ProductGet failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductPut(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	tracer := otel.Tracer("product.controller")
	ctx, span := tracer.Start(ctx, "ProductController.ProductPut")
	defer span.End()
	
	logger.Info(ctx, "product controller ProductPut called")

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
		logger.Error(ctx, "product controller ProductPut failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}