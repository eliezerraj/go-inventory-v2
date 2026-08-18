package repository

import (
	"context"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/go-inventory-v2/application/domain/entity"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type InventoryPriceRepository struct {
	dbConnector connector.IDatabaseConnector
}

type IInventoryPriceRepository interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	InventoryAdd(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error)
	InventoryGet(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error)
	InventoryPut(ctx context.Context, inventory entity.Inventory) (int64, error)
	PriceAdd(ctx context.Context, price entity.Price) (*entity.Price, error)
	PriceGet(ctx context.Context, price entity.Price) (*entity.Price, error)
	PricePut(ctx context.Context, price entity.Price) (int64, error)
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

func (ipr *InventoryPriceRepository) InventoryAdd(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryAdd")
	defer span.End()

	logger.Info(ctx, "inventory price repository InventoryAdd called")

	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryAdd failed", zap.Error(err))
		}
	}()

	connectorWriter := ipr.dbConnector.Writer()

	query := `INSERT INTO inventory ( fk_product_id, 
										available,
										pending,
										sold,
										created_at) 
				VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id int
	if err := connectorWriter.QueryRow(ctx, query, inventory.ProductId, inventory.Available, inventory.Pending, inventory.Sold, inventory.CreatedAt).Scan(&id); err != nil {
		return nil, err
	}
	inventory.ID = id

	return &inventory, nil
}

func (ipr *InventoryPriceRepository) InventoryGet(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryGet")
	defer span.End()

	logger.Info(ctx, "inventory price repository InventoryGet called")
	
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryGet failed", zap.Error(err))
		}
	}()

	connectorReader := ipr.dbConnector.Reader()

	query := `SELECT id, fk_product_id, available, pending, sold, created_at, updated_at 
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

func (ipr *InventoryPriceRepository) InventoryPut(ctx context.Context, inventory entity.Inventory) (int64, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.InventoryPut")
	defer span.End()

	logger.Info(ctx, "inventory price repository InventoryPut called", zap.Any("inventory", inventory))

	connectorWriter := ipr.dbConnector.Writer()
	
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository InventoryPut failed", zap.Error(err))
		}
	}()

	query := `UPDATE inventory 
				SET available = $1,
					pending = $2,
					sold = $3,
					updated_at = $4
				WHERE fk_product_id = $5`

	row, err := connectorWriter.Exec(ctx, query, inventory.Available, inventory.Pending, inventory.Sold, inventory.UpdatedAt, inventory.ProductId)
	if err != nil {
		return 0, err
	}

	rowsAffected := row.RowsAffected()
	return rowsAffected, nil
}

func (ipr *InventoryPriceRepository) PriceGet(ctx context.Context, price entity.Price) (*entity.Price, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PriceGet")
	defer span.End()

	logger.Info(ctx, "inventory price repository PriceGet called")

	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PriceGet failed", zap.Error(err))
		}
	}()
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

func (ipr *InventoryPriceRepository) PriceAdd(ctx context.Context, price entity.Price) (*entity.Price, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PriceAdd")
	defer span.End()
	
	logger.Info(ctx, "inventory price repository PriceAdd called")

	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PriceAdd failed", zap.Error(err))
		}
	}()

	connectorWriter := ipr.dbConnector.Writer()

	query := `INSERT INTO price ( fk_product_id,
									amount,
									currency,
									started_at,
									ended_at,
									created_at) 
				VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	if err := connectorWriter.QueryRow(ctx, query, price.ProductId, price.Amount, price.Currency, price.StartedAt, price.EndedAt, price.CreatedAt).Scan(&id); err != nil {
		return nil, err
	}
	price.ID = id

	return &price, nil
}

func (ipr *InventoryPriceRepository) PricePut(ctx context.Context, price entity.Price) (int64, error) {
	tracer := otel.Tracer("inventory_price.repository")
	ctx, span := tracer.Start(ctx, "InventoryPriceRepository.PricePut")
	defer span.End()

	logger.Info(ctx, "inventory price repository PricePut called", zap.Any("price", price))

	var err error

	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "inventory price repository PricePut failed", zap.Error(err))
		}
	}()

	connectorWriter := ipr.dbConnector.Writer()

	query := `UPDATE price 
				SET amount = $1,
					currency = $2,
					started_at = $3,
					ended_at = $4,
					updated_at = $5
				WHERE fk_product_id = $6`

	row, err := connectorWriter.Exec(ctx, query, price.Amount, price.Currency, price.StartedAt, price.EndedAt, price.UpdatedAt, price.ProductId)
	if err != nil {
		return 0, err
	}

	rowsAffected := row.RowsAffected()
	return rowsAffected, nil
}