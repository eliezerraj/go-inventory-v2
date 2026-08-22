package repository

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/go-inventory-v2/application/domain/entity"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type InventoryPriceRepository struct {
	dbConnector connector.IDatabaseConnector
}

type IInventoryPriceRepository interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	InventoryAdd(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (*entity.Inventory, error)
	InventoryGet(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error)
	InventoryPut(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (int64, error)
	InventoryPatch(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (int64, error)
	PriceAdd(ctx context.Context, tx pgx.Tx, price entity.Price) (*entity.Price, error)
	PriceGet(ctx context.Context, price entity.Price) (*entity.Price, error)
	PricePut(ctx context.Context, tx pgx.Tx, price entity.Price) (int64, error)
}

func NewInventoryPriceRepository(dbConnector connector.IDatabaseConnector) IInventoryPriceRepository {
	logger.InfoOutCtx("initializing inventory price repository SUCCESSFULLY")
	return &InventoryPriceRepository{
		dbConnector: dbConnector,
	}
}

func (ipr *InventoryPriceRepository) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	logger.Info(ctx, "inventory price repository BeginTx called")

	tx, err := ipr.dbConnector.Writer().BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ipr *InventoryPriceRepository) InventoryAdd(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (res_inventory *entity.Inventory, err error) {
	logger.Info(ctx, "inventory price repository InventoryAdd called")

	// Start tracing and metrics
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryAdd")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
    counter, _ := meter.Int64Counter("db_custom_inventory_add_requests_total")
    histogram, _ := meter.Float64Histogram("db_custom_inventory_add_duration_seconds")
    start := time.Now()

    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "InventoryAdd"),
    ))
	
	// Defer function to handle error logging and metrics recording
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryAdd failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "InventoryAdd"),
        ))
	}()

	// Insert inventory record into the database
	query := `INSERT INTO inventory ( fk_product_id, 
										available,
										pending,
										sold,
										created_at) 
				VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id int
	if err := tx.QueryRow(ctx, query, inventory.ProductId, inventory.Available, inventory.Pending, inventory.Sold, inventory.CreatedAt).Scan(&id); err != nil {
		return nil, err
	}
	inventory.ID = id

	return &inventory, nil
}

func (ipr *InventoryPriceRepository) InventoryGet(ctx context.Context, inventory entity.Inventory) (res_inventory *entity.Inventory, err error) {
	logger.Info(ctx, "inventory price repository InventoryGet called")

	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryGet")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
    counter, _ := meter.Int64Counter("db_custom_inventory_get_requests_total")
    histogram, _ := meter.Float64Histogram("db_custom_inventory_get_duration_seconds")
    start := time.Now()

    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "InventoryGet"),
    ))
	
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryGet failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "InventoryGet"),
        ))
	}()

	connectorReader := ipr.dbConnector.Reader()

	query := `SELECT id, 
					fk_product_id, 
					available, 
					pending, 
					sold, 
					created_at, 
					updated_at 
				FROM inventory 
				WHERE fk_product_id = $1`

	var inv entity.Inventory
	err = connectorReader.QueryRow(ctx, query, inventory.ProductId).Scan(&inv.ID, &inv.ProductId, &inv.Available, &inv.Pending, &inv.Sold, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		logger.Warn(ctx, "failed to get inventory record", zap.Error(err))
		return nil, err
	}

	return &inv, nil
}

func (ipr *InventoryPriceRepository) InventoryPut(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (rowsAffected int64, err error) {
	logger.Info(ctx, "inventory price repository InventoryPut called", zap.Any("inventory", inventory))

	// Start tracing and metrics
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryPut")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_inventory_put_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_inventory_put_duration_seconds")
	start := time.Now()

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "InventoryPut"),
	))

	// Defer function to handle error logging and metrics recording
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryPut failed", zap.Error(err))
		}
		histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("operation", "InventoryPut"),
		))
	}()

	// Update inventory record in the database
	query := `UPDATE inventory 
				SET available = $1,
					pending = $2,
					sold = $3,
					updated_at = $4
				WHERE fk_product_id = $5`

	row, err := tx.Exec(ctx, query, inventory.Available, inventory.Pending, inventory.Sold, inventory.UpdatedAt, inventory.ProductId)
	if err != nil {
		return 0, err
	}

	rowsAffected = row.RowsAffected()
	return rowsAffected, nil
}

func (ipr *InventoryPriceRepository) InventoryPatch(ctx context.Context, tx pgx.Tx, inventory entity.Inventory) (rowsAffected int64, err error) {
	logger.Info(ctx, "inventory price repository InventoryPatch called")

	tracer := otel.Tracer("inventory_price.repository")
    ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryPatch")
    defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_inventory_patch_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_inventory_patch_duration_seconds")
	start := time.Now()
	
	counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "InventoryPatch"),
    ))
	
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryPatch failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "InventoryPatch"),
        ))
	}()

	query := `UPDATE inventory 
				SET available = available + $1,
					sold = sold + $2,
					pending = pending + $3,
					updated_at = $4
				WHERE id = $5`

	row, err := tx.Exec(ctx, query, inventory.Available, inventory.Sold, inventory.Pending, inventory.UpdatedAt, inventory.ID)
	if err != nil {
		return 0, err
	}

	if row.RowsAffected() == 0 {
		logger.Warn(ctx, "inventory price repository InventoryPatch: no rows affected, inventory not found", zap.Int("inventory_id", inventory.ID))
		return 0, nil
	}

	rowsAffected = row.RowsAffected()
	return rowsAffected, nil
}

func (ipr *InventoryPriceRepository) PriceGet(ctx context.Context, price entity.Price) (res_price *entity.Price, err error) {
	logger.Info(ctx, "inventory price repository PriceGet called")

	// Start tracing and metrics
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PriceGet")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_price_get_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_price_get_duration_seconds")
	start := time.Now()

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "PriceGet"),
	))

	// Defer function to handle error logging and metrics recording
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PriceGet failed", zap.Error(err))
			logger.Error(ctx, "inventory price repository PriceGet failed", zap.Error(err))
		}
		histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("operation", "PriceGet"),
		))
	}()

	// Query the database for the price record
	connectorReader := ipr.dbConnector.Reader()

	query := `SELECT id, fk_product_id, amount, currency, started_at, ended_at, created_at, updated_at 
				FROM price 
				WHERE fk_product_id = $1`

	var pr entity.Price
	err = connectorReader.QueryRow(ctx, query, price.ProductId).Scan(&pr.ID, &pr.ProductId, &pr.Amount, &pr.Currency, &pr.StartedAt, &pr.EndedAt, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		logger.Warn(ctx, "failed to get price record", zap.Error(err))
		return nil, err
	}

	return &pr, nil
}

func (ipr *InventoryPriceRepository) PriceAdd(ctx context.Context, tx pgx.Tx, price entity.Price) (res_price *entity.Price, err error) {
	logger.Info(ctx, "inventory price repository PriceAdd called")

	// Start tracing and metrics
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PriceAdd")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_price_add_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_price_add_duration_seconds")
	start := time.Now()

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "PriceAdd"),
	))

	// Defer function to handle error logging and metrics recording
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PriceAdd failed", zap.Error(err))
		}
		histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("operation", "PriceAdd"),
		))
	}()

	// Insert price record into the database
	query := `INSERT INTO price ( fk_product_id,
									amount,
									currency,
									started_at,
									ended_at,
									created_at) 
				VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	if err := tx.QueryRow(ctx, query, price.ProductId, price.Amount, price.Currency, price.StartedAt, price.EndedAt, price.CreatedAt).Scan(&id); err != nil {
		return nil, err
	}

	price.ID = id
	return &price, nil
}

func (ipr *InventoryPriceRepository) PricePut(ctx context.Context, tx pgx.Tx, price entity.Price) (rowsAffected int64, err error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PricePut")
	defer span.End()

	meter := otel.Meter("go-inventory-v2.repository")
	counter, _ := meter.Int64Counter("db_custom_price_put_requests_total")
	histogram, _ := meter.Float64Histogram("db_custom_price_put_duration_seconds")
	start := time.Now()

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "PricePut"),
	))

	logger.Info(ctx, "inventory price repository PricePut called", zap.Any("price", price))

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PricePut failed", zap.Error(err))
		}
		histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("operation", "PricePut"),
		))
	}()

	query := `UPDATE price 
				SET amount = $1,
					currency = $2,
					started_at = $3,
					ended_at = $4,
					updated_at = $5
				WHERE fk_product_id = $6`

	row, err := tx.Exec(ctx, query, price.Amount, price.Currency, price.StartedAt, price.EndedAt, price.UpdatedAt, price.ProductId)
	if err != nil {
		return 0, err
	}

	rowsAffected = row.RowsAffected()
	return rowsAffected, nil
}