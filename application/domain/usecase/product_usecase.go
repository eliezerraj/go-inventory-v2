package usecase

import (
	"time"
	"context"

	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/domain/entity"
	"github.com/go-inventory-v2/application/infrastructure/repository"

	"github.com/jackc/pgx/v5"
)

type ProductUsecase struct {
	productRepository repository.IProductRepository
}

type IProductUseCase interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
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

func (p *ProductUsecase) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	logger.InfoOutCtx("product usecase BeginTx called")

	tx, err := p.productRepository.BeginTx(ctx, opts)
	if err != nil {
		logger.ErrorOutCtx("product usecase BeginTx failed", zap.Error(err))
		return nil, err
	}

	return tx, nil
}

func (p *ProductUsecase) ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.InfoOutCtx("product usecase ProductAdd called")

	tx, err := p.productRepository.BeginTx(ctx, pgx.TxOptions{ IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite })
	if err != nil {
		logger.ErrorOutCtx("product usecase ProductAdd failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.ErrorOutCtx("product usecase ProductAdd failed to rollback transaction", zap.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				logger.ErrorOutCtx("product usecase ProductAdd failed to commit transaction", zap.Error(commitErr))
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
		logger.ErrorOutCtx("product usecase ProductAdd failed", zap.Error(err))
		return nil, err
	}

	if err != nil {
		logger.ErrorOutCtx("product usecase ProductAdd failed to commit transaction", zap.Error(err))
		return nil, err
	}

	logger.InfoOutCtx("product usecase ProductAdd completed SUCCESSFULLY")
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