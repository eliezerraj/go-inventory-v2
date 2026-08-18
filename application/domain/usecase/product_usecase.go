package usecase

import (
	"time"
	"errors"
	"context"

	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/domain/entity"
	"github.com/go-inventory-v2/application/infrastructure/repository"

	"github.com/jackc/pgx/v5"

	"go.opentelemetry.io/otel"
)

type ProductUsecase struct {
	productRepository repository.IProductRepository
	inventoryPriceRepository repository.IInventoryPriceRepository
}

type IProductUseCase interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductPut(ctx context.Context, product entity.Product) (*entity.Product, error)
}

func NewProductUseCase(productRepository repository.IProductRepository, 
						inventoryPriceRepository repository.IInventoryPriceRepository) IProductUseCase {
	logger.InfoOutCtx("initializing product usecase SUCCESSFULLY")

	return &ProductUsecase{
		productRepository: productRepository,
		inventoryPriceRepository: inventoryPriceRepository,
	}
}

func (p *ProductUsecase) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	logger.Info(ctx, "product usecase BeginTx called")

	tx, err := p.productRepository.BeginTx(ctx, opts)
	if err != nil {
		logger.Error(ctx, "product usecase BeginTx failed", zap.Error(err))
		return nil, err
	}

	return tx, nil
}

func (p *ProductUsecase) ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error) {
	tracer := otel.Tracer("inventory.repository")
	ctx, span := tracer.Start(ctx, "ProductUsecase.ProductAdd")
	defer span.End()

	logger.Info(ctx, "product usecase ProductAdd called")

	tx, err := p.productRepository.BeginTx(ctx, pgx.TxOptions{ IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite })
	if err != nil {
		logger.Error(ctx, "product usecase ProductAdd failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.Error(ctx, "product usecase ProductAdd failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				logger.Error(ctx, "product usecase ProductAdd failed to commit transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	expiresAt := time.Now().AddDate(0, 0, 360).UTC() // Set expires_at to 360 days from now
	product.ExpiresAt = &expiresAt
	createAt := time.Now().UTC()
	product.CreatedAt = &createAt 

	res, err := p.productRepository.ProductAdd(ctx, product)
	if err != nil {
		logger.Error(ctx, "product usecase ProductAdd failed", zap.Error(err))
		return nil, err
	}

	if err != nil {
		logger.Error(ctx, "product usecase ProductAdd failed to commit transaction", zap.Error(err))
		return nil, err
	}

	if product.Inventory != nil {
		product.Inventory.ProductId = res.ID
		product.Inventory.CreatedAt = &createAt
		inv, err := p.inventoryPriceRepository.InventoryAdd(ctx, *product.Inventory)
		if err != nil {
			logger.Error(ctx, "product usecase ProductAdd failed to add inventory", zap.Error(err))
			return nil, err
		}
		res.Inventory = inv
	}

	if product.Price != nil {
		product.Price.ProductId = res.ID
		product.Price.CreatedAt = &createAt
		product.Price.StartedAt = &createAt
		product.Price.EndedAt = &expiresAt

		price, err := p.inventoryPriceRepository.PriceAdd(ctx, *product.Price)
		if err != nil {
			logger.Error(ctx, "product usecase ProductAdd failed to add price", zap.Error(err))
			return nil, err
		}

		res.Price = price
	}

	logger.Info(ctx, "product usecase ProductAdd completed SUCCESSFULLY")
	return res, nil
}

func (p *ProductUsecase) ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error) {
	tracer := otel.Tracer("inventory.repository")
	ctx, span := tracer.Start(ctx, "ProductUsecase.ProductGet")
	defer span.End()

	logger.Info(ctx, "product usecase ProductGet called")

	res, err := p.productRepository.ProductGet(ctx, product)
	if err != nil {
		logger.Error(ctx, "product usecase ProductGet failed", zap.Error(err))
		return nil, err
	}

	inv, err := p.inventoryPriceRepository.InventoryGet(ctx, entity.Inventory{ProductId: res.ID})
	if err != nil {
		logger.Warn(ctx, "product usecase ProductGet failed to get inventory", zap.Error(err))
	}
	res.Inventory = inv

	price, err := p.inventoryPriceRepository.PriceGet(ctx, entity.Price{ProductId: res.ID})
	if err != nil {
		logger.Warn(ctx, "product usecase ProductGet failed to get price", zap.Error(err))
	}
	res.Price = price

	return res, nil
}

func (p *ProductUsecase) ProductPut(ctx context.Context, product entity.Product) (*entity.Product, error) {
	tracer := otel.Tracer("inventory.repository")
	ctx, span := tracer.Start(ctx, "ProductUsecase.ProductPut")
	defer span.End()

	logger.Info(ctx, "product usecase ProductPut called")

	// Start tx
	tx, err := p.productRepository.BeginTx(ctx, pgx.TxOptions{ IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite })
	if err != nil {
		logger.Error(ctx, "product usecase ProductPut failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.Error(ctx, "product usecase ProductPut failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				logger.Error(ctx, "product usecase ProductPut failed to commit transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	// Check if the product exists before updating
	res_prod, err := p.productRepository.ProductGet(ctx, product)
	if err != nil {
		logger.Error(ctx, "product usecase ProductPut failed", zap.Error(err))
		tx.Rollback(ctx)
		return nil, err
	}

	product.ID = res_prod.ID
	updatedAt := time.Now().UTC()
	product.UpdatedAt = &updatedAt

	// Update product
	upd_prod, err := p.productRepository.ProductPut(ctx, product)
	if err != nil {
		logger.Error(ctx, "product usecase ProductPut failed", zap.Error(err))
		return nil, err
	}

	if upd_prod == 0{
		logger.Warn(ctx, "product usecase ProductPut: no rows affected, product not found", zap.Int("product_id", product.ID))
		return nil, errors.New("product not found")
	}

	// Update price
	res_price, err := p.inventoryPriceRepository.PriceGet(ctx, entity.Price{ProductId: res_prod.ID})
	if err != nil {
		logger.Warn(ctx, "product usecase ProductPut failed to get inventory", zap.Error(err))
		return nil, err
	}

	product.Price.ProductId = product.ID
	product.Price.UpdatedAt = &updatedAt

	if product.Price.StartedAt == nil{
		product.Price.StartedAt = res_price.StartedAt
	}

	if product.Price.EndedAt == nil{
		product.Price.EndedAt = res_price.EndedAt
	}

	upd_price, err := p.inventoryPriceRepository.PricePut(ctx, *product.Price)
	if err != nil {
		logger.Error(ctx, "product usecase ProductPut failed", zap.Error(err))
		return nil, err
	}

	if upd_price == 0{
		logger.Warn(ctx, "product usecase ProductPut: no rows affected, price not found", zap.Int("product_id", product.ID))
		return nil, errors.New("price not found")
	}

	// Update inventory
	product.Inventory.ProductId = product.ID
	product.Inventory.UpdatedAt = &updatedAt
	upd_inv, err := p.inventoryPriceRepository.InventoryPut(ctx, *product.Inventory)
	if err != nil {
		logger.Error(ctx, "product usecase ProductPut failed", zap.Error(err))
		return nil, err
	}

	if upd_inv == 0{
		logger.Warn(ctx, "product usecase ProductPut: no rows affected, inventory not found", zap.Int("product_id", product.ID))
		return nil, errors.New("inventory not found")
	}
	
	return &product, nil
}