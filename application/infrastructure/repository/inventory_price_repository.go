package repository

import (
	"context"
	//"errors"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/go-inventory-v2/application/domain/entity"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"
)

type InventoryPriceRepository struct {
	dbConnector connector.IDatabaseConnector
}

type IInventoryPriceRepository interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	InventoryAdd(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error)
	InventoryGet(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error)
	PriceAdd(ctx context.Context, price entity.Price) (*entity.Price, error)
	PriceGet(ctx context.Context, price entity.Price) (*entity.Price, error)
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
	logger.Info(ctx, "inventory price repository InventoryAdd called")

	connectorWriter := ipr.dbConnector.Writer()

	query := `INSERT INTO inventory ( fk_product_id, 
										available,
										pending,
										sold,
										created_at) 
				VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id int
	if err := connectorWriter.QueryRow(ctx, query, inventory.ProductId, inventory.Available, inventory.Pending, inventory.Sold, inventory.CreatedAt).Scan(&id); err != nil {
		logger.Error(ctx, "failed to insert inventory record", zap.Error(err))
		return nil, err
	}
	inventory.ID = id

	return &inventory, nil
}

func (ipr *InventoryPriceRepository) InventoryGet(ctx context.Context, inventory entity.Inventory) (*entity.Inventory, error) {
	logger.Info(ctx, "inventory price repository InventoryGet called")

	connectorReader := ipr.dbConnector.Reader()

	query := `SELECT id, fk_product_id, available, pending, sold, created_at, updated_at 
				FROM inventory 
				WHERE fk_product_id = $1`

	var inv entity.Inventory
	err := connectorReader.QueryRow(ctx, query, inventory.ProductId).Scan(&inv.ID, &inv.ProductId, &inv.Available, &inv.Pending, &inv.Sold, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		logger.Warn(ctx, "failed to get inventory record", zap.Error(err))
		return nil, err
	}

	return &inv, nil
}

func (ipr *InventoryPriceRepository) PriceGet(ctx context.Context, price entity.Price) (*entity.Price, error) {
	logger.Info(ctx, "inventory price repository PriceGet called")

	connectorReader := ipr.dbConnector.Reader()

	query := `SELECT id, fk_product_id, amount, currency, started_at, ended_at, created_at, updated_at 
				FROM price 
				WHERE fk_product_id = $1`

	var pr entity.Price
	err := connectorReader.QueryRow(ctx, query, price.ProductId).Scan(&pr.ID, &pr.ProductId, &pr.Amount, &pr.Currency, &pr.StartedAt, &pr.EndedAt, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		logger.Warn(ctx, "failed to get price record", zap.Error(err))
		return nil, err
	}

	return &pr, nil
}

func (ipr *InventoryPriceRepository) PriceAdd(ctx context.Context, price entity.Price) (*entity.Price, error) {
	logger.Info(ctx, "inventory price repository PriceAdd called")

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
		logger.Error(ctx, "failed to insert price record", zap.Error(err))
		return nil, err
	}
	price.ID = id

	return &price, nil
}