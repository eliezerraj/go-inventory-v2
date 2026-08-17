package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/go-inventory-v2/application/domain/entity"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"
)

type ProductRepository struct {
	dbConnector connector.IDatabaseConnector
}

type IProductRepository interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductPut(ctx context.Context, product entity.Product) (int64, error)
}

func NewProductRepository(dbConnector connector.IDatabaseConnector) IProductRepository {
	logger.InfoOutCtx("initializing product repository SUCCESSFULLY")

	return &ProductRepository{
		dbConnector: dbConnector,
	}
}

func (p *ProductRepository) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	logger.InfoOutCtx("product repository BeginTx called")

	tx, err := p.dbConnector.Writer().BeginTx(ctx, opts)
	if err != nil {
		logger.ErrorOutCtx("product repository BeginTx failed", zap.Error(err))
		return nil, err
	}

	return tx, nil
}

func (p *ProductRepository) ProductAdd(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.Info(ctx, "product repository ProductAdd called")

	connectorWriter := p.dbConnector.Writer()

	query := `INSERT INTO product ( sku, 
									type,
									name,
									status,
									lead_time,
									expires_at,
									created_at) 
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	rows := connectorWriter.QueryRow(ctx, query, product.Sku, product.Type, product.Name, product.Status, product.LeadTime, product.ExpiresAt, product.CreatedAt)

	var id int
	if err := rows.Scan(&id); err != nil {
		logger.ErrorOutCtx("product repository ProductAdd failed", zap.Error(err))
		return nil, err
	}

	product.ID = id
	return &product, nil
}

func (p *ProductRepository) ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error) {
	logger.Info(ctx, "product repository ProductGet called")

	// Get a reader connection from the database connector
	connectorReader := p.dbConnector.Reader()

	query := `SELECT id,
					 sku, 
					 name, 
					 status, 
					 type,
					 lead_time,
					 expires_at,
					 created_at,
					 updated_at
	 		  FROM product WHERE sku = $1`

	rows, err := connectorReader.Query(ctx, query, product.Sku)
	if err != nil {
		logger.Error(ctx,"failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Status, &product.Type, &product.LeadTime, &product.ExpiresAt, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			logger.Error(ctx, "failed to scan result", zap.Error(err))
			return nil, err
		}
	} else {
		logger.Warn(ctx, "not found", zap.String("sku", product.Sku))
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func (p *ProductRepository) ProductPut(ctx context.Context, product entity.Product) (int64, error) {
	logger.Info(ctx, "product repository ProductPut called")

	connectorWriter := p.dbConnector.Writer()

	logger.Info(ctx, "product repository ProductPut called", zap.Any("product", product))

	query := `UPDATE product 
				SET sku = $1,
					type = $2,
					name = $3,
					status = $4,
					lead_time = $5,
					expires_at = $6,
					updated_at = $7
				WHERE id = $8`

	row, err := connectorWriter.Exec(ctx, query, product.Sku, product.Type, product.Name, product.Status, product.LeadTime, product.ExpiresAt, product.UpdatedAt, product.ID)
	if err != nil {
		logger.ErrorOutCtx("product repository ProductPut failed", zap.Error(err))
		return 0, err
	}

	if row.RowsAffected() == 0 {
		logger.Warn(ctx, "product repository ProductPut: no rows affected, product not found", zap.Int("product_id", product.ID))
		return 0, nil
	}

	return row.RowsAffected(), nil
}