package repository

import (
	"time"
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/go-inventory-v2/application/domain/entity"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/attribute"
)

type ProductRepository struct {
	dbConnector connector.IDatabaseConnector
}

type IProductRepository interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	ProductAdd(ctx context.Context, tx pgx.Tx, product entity.Product) (*entity.Product, error)
	ProductGet(ctx context.Context, product entity.Product) (*entity.Product, error)
	ProductPut(ctx context.Context, tx pgx.Tx, product entity.Product) (int64, error)
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

func (p *ProductRepository) ProductAdd(ctx context.Context, tx pgx.Tx, product entity.Product) (res_product *entity.Product, err error) {
	logger.Info(ctx, "product repository ProductAdd called")

	tracer := otel.Tracer("inventory.repository")
	ctx, span := tracer.Start(ctx, "ProductRepository.ProductAdd")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_product_add_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_product_add_duration_seconds")
	start := time.Now()

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "ProductAdd"),
	))

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "product repository ProductAdd failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "ProductAdd"),
        ))
	}()

	query := `INSERT INTO product ( sku, 
									type,
									name,
									status,
									lead_time,
									expires_at,
									created_at) 
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	rows := tx.QueryRow(ctx, query, product.Sku, product.Type, product.Name, product.Status, product.LeadTime, product.ExpiresAt, product.CreatedAt)

	var id int
	if err = rows.Scan(&id); err != nil {
		return nil, err
	}

	product.ID = id
	return &product, nil
}

func (p *ProductRepository) ProductGet(ctx context.Context, product entity.Product) (res_product *entity.Product, err error) {
	logger.Info(ctx, "product repository ProductGet called")

	tracer := otel.Tracer("inventory.repository")
    ctx, span := tracer.Start(ctx, "ProductRepository.ProductGet")
    defer span.End()

    meter := otel.Meter("go-inventory-v2.repository")
    counter, _ := meter.Int64Counter("db_custom_product_get_requests_total")
    histogram, _ := meter.Float64Histogram("db_custom_product_get_duration_seconds")
    start := time.Now()

    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "ProductGet"),
    ))

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "product repository ProductGet failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "ProductGet"),
        ))
	}()

	// Get a reader connection from the database connector
	connectorReader := p.dbConnector.Reader()
	
	var where string
	if product.ID != 0 {
		where = "id = $1"
	} else {
		where = "sku = $1"
	}

	query := `SELECT id,
					 sku, 
					 name, 
					 status, 
					 type,
					 lead_time,
					 expires_at,
					 created_at,
					 updated_at
	 		  FROM product WHERE ` + where

	var arg interface{}
	if product.Sku == "" {
		arg = product.ID
	} else {
		arg = product.Sku
	}

	var rows pgx.Rows
	rows, err = connectorReader.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Status, &product.Type, &product.LeadTime, &product.ExpiresAt, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res_product = &product
	} else {
		logger.Warn(ctx, "not found", zap.String("sku", product.Sku))
		err = errors.New("product not found")
		return nil, err
	}

	return res_product, nil
}

func (p *ProductRepository) ProductPut(ctx context.Context, tx pgx.Tx, product entity.Product) (rowsAffected int64, err error) {
	logger.Info(ctx, "product repository ProductPut called")

	tracer := otel.Tracer("inventory.repository")
    ctx, span := tracer.Start(ctx, "ProductRepository.ProductPut")
    defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_product_put_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_product_put_duration_seconds")
	start := time.Now()
	
	counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "ProductPut"),
    ))
	
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "product repository ProductPut failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "ProductPut"),
        ))
	}()

	query := `UPDATE product 
				SET sku = $1,
					type = $2,
					name = $3,
					status = $4,
					lead_time = $5,
					expires_at = $6,
					updated_at = $7
				WHERE id = $8`

	row, err := tx.Exec(ctx, query, product.Sku, product.Type, product.Name, product.Status, product.LeadTime, product.ExpiresAt, product.UpdatedAt, product.ID)
	if err != nil {
		return 0, err
	}

	if row.RowsAffected() == 0 {
		logger.Warn(ctx, "product repository ProductPut: no rows affected, product not found", zap.Int("product_id", product.ID))
		return 0, nil
	}

	rowsAffected = row.RowsAffected()
	return rowsAffected, nil
}
