package controller

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/domain/usecase"
	"github.com/go-inventory-v2/application/domain/external"
	"github.com/go-inventory-v2/application/domain/entity"
	"github.com/go-inventory-v2/application/tracing"

	"go.opentelemetry.io/otel/trace"
)

type ProductController struct {
	schema Schema
	productUseCase usecase.IProductUseCase
}

// Schema struct defines a validation schema for product requests.
type Schema struct {
    Validate func(context.Context, any) error
}

func inventoryPatchSchema() Schema {
    return Schema{
        Validate: func(ctx context.Context, data any) error {
			
			req, ok := data.(external.ProductRequest)
			if !ok {
                return errors.New("schema validation failed ! Please check the request body and try again.")
            }

			if req.Inventory == nil {
                return errors.New("schema validation failed ! field inventory is mandatory")
            }

            return nil
        },
    }
}

// Use in ProductAdd and ProductPut.
func productAddSchema() Schema {
    return Schema{
        Validate: func(ctx context.Context, data any) error {
            
			req, ok := data.(external.ProductRequest)
			if !ok {
                return errors.New("schema validation failed ! Please check the request body and try again.")
            }

			validate := validator.New()
			err := validate.Struct(req)
			if err != nil {
				return errors.New("schema validation failed ! Please check the request body data and try again.")
			}

            if req.Price == nil {
                return errors.New("schema validation failed ! field price is mandatory")
            }

            if req.Inventory == nil {
                return errors.New("schema validation failed ! field inventory is mandatory")
            }

            return nil
        },
    }
}

// NewProductController creates a new instance of ProductController with the provided product use case.
func NewProductController(productUseCase usecase.IProductUseCase) *ProductController {
	logger.InfoOutCtx("initializing product controller SUCCESSFULLY")

	schema := Schema{
		Validate: func(ctx context.Context, data any) error {
			return nil
		},
	}
	return &ProductController{
		schema:           schema,
		productUseCase:   productUseCase,
	}
}

func (p *ProductController) ProductAdd(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.Info(ctx, "product controller ProductAdd called")

	// Tracing and metrics
	ctx, span := tracing.CustomStartSpanCtx(ctx, "productController.productAdd", trace.SpanKindInternal)
	defer span.End()

	// Schema validation
	schema := productAddSchema()
    if err := schema.Validate(ctx, req); err != nil {
        return nil, err
    }

	// Create a Product entity from the request
	product := entity.Product{
		Sku:         req.Sku,
		Name:        req.Name,
		Status:      req.Status,
		Type:        req.Type,
		LeadTime:    req.LeadTime,
	}

	price := entity.Price{
		Currency:    req.Price.Currency,
		Amount:      req.Price.Amount,
	}
	product.Price = &price

	inventory := entity.Inventory{
		Available:   req.Inventory.Available,
		Pending:     req.Inventory.Pending,
		Sold:        req.Inventory.Sold,
	}
	product.Inventory = &inventory
	
	// Call the use case to add the product
	res, err := p.productUseCase.ProductAdd(ctx, product)
	if err != nil {
		logger.Error(ctx, "product controller ProductAdd failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductGet(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.Info(ctx, "product controller ProductGet called", zap.String("sku", req.Sku))
	
	// Tracing and metrics
	ctx, span := tracing.CustomStartSpanCtx(ctx, "productController.productGet", trace.SpanKindInternal)
	defer span.End()

	product := entity.Product{
		ID:  req.ID,
		Sku: req.Sku,
	}

	// Call the use case to get the product
	res, err := p.productUseCase.ProductGet(ctx, product)
	if err != nil {
		logger.Error(ctx, "product controller ProductGet failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) ProductPut(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.Info(ctx, "product controller ProductPut called")

	// Tracing and metrics
	ctx, span := tracing.CustomStartSpanCtx(ctx, "productController.productPut", trace.SpanKindInternal)
	defer span.End()
	
	// Schema validation
	schema := productAddSchema()
	if err := schema.Validate(ctx, req); err != nil {
		return nil, err
	}

	product := entity.Product{
		Sku:         req.Sku,
		Name:        req.Name,
		Status:      req.Status,
		Type:        req.Type,
		LeadTime:    req.LeadTime,
	}

	price := entity.Price{
		Currency:    req.Price.Currency,
		Amount:      req.Price.Amount,
	}
	product.Price = &price

	inventory := entity.Inventory{
		Available:   req.Inventory.Available,
		Pending:     req.Inventory.Pending,
		Sold:        req.Inventory.Sold,
	}
	
	product.Inventory = &inventory

	// Call the use case to update the product
	res, err := p.productUseCase.ProductPut(ctx, product)
	if err != nil {
		logger.Error(ctx, "product controller ProductPut failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (p *ProductController) InventoryPatch(ctx context.Context, req external.ProductRequest) (*entity.Product, error) {
	logger.Info(ctx, "product controller InventoryPatch called")

	// Tracing and metrics
	ctx, span := tracing.CustomStartSpanCtx(ctx, "productController.inventoryPatch", trace.SpanKindInternal)
	defer span.End()
	
	// Schema validation
	schema := inventoryPatchSchema()
	if err := schema.Validate(ctx, req); err != nil {
		return nil, err
	}

	product := entity.Product{
		ID:  req.ID,
		Sku: req.Sku,
	}

	inventory := entity.Inventory{
		Available:   req.Inventory.Available,
		Pending:     req.Inventory.Pending,
		Sold:        req.Inventory.Sold,
	}

	product.Inventory = &inventory

	// Call the use case to patch the inventory
	res, err := p.productUseCase.InventoryPatch(ctx, product)
	if err != nil {
		logger.Error(ctx, "product controller InventoryPatch failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}