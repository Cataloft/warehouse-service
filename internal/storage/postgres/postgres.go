package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Cataloft/warehouse-service/internal/config"
	errs "github.com/Cataloft/warehouse-service/internal/constants/errors"
	"github.com/Cataloft/warehouse-service/internal/handlers/warehouses"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Conn   *pgxpool.Pool
	logger *slog.Logger
}

func New(cfg *config.Database, logger *slog.Logger) *Storage {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error: parse db url")
	}

	var connPool *pgxpool.Pool

	for i := cfg.MaxAttempts; i > 0; i-- {
		connPool, _ = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if connPool.Ping(context.Background()) == nil {
			logger.Debug("postgres upped")

			break
		}

		time.Sleep(cfg.AttemptDuration)
	}

	return &Storage{
		Conn:   connPool,
		logger: logger,
	}
}

func (s *Storage) UpdateGood(uniqueCode string, amount int32) (int32, error) {
	query := "UPDATE goods SET amount = amount - $1 WHERE unique_code = $2 RETURNING amount"

	tx, err := s.Conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = tx.Rollback(context.Background()); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("Error transaction Rollback", "error", err)
		}
	}()

	row := tx.QueryRow(context.Background(), query, amount, uniqueCode)

	err = row.Scan(&amount)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (s *Storage) GetAvailableGood(uniqueCode string, amount int32) error {
	query :=
		`SELECT amount, availability 
			FROM goods 
			JOIN warehouses on goods.warehouse_id = warehouses.id 
			WHERE unique_code = $1 AND availability = true`

	rows, err := s.Conn.Query(context.Background(), query, uniqueCode)
	if err != nil {
		return err
	}

	var (
		availability  bool
		currentAmount int32
	)

	for rows.Next() {
		if err := rows.Scan(&currentAmount, &availability); err != nil {
			return err
		}

		if (currentAmount - amount) < 0 {
			return fmt.Errorf("%s", errs.GoodError)
		}
	}

	return nil
}

func (s *Storage) GetWarehouse(warehouseID int64) (warehouses.Attributes, int32, error) {
	query :=
		`SELECT warehouses.name, availability, SUM(goods.amount) as total 
			FROM warehouses 
			JOIN goods on warehouses.id = goods.warehouse_id 
			WHERE warehouses.id = $1 
			GROUP BY warehouses.name, availability`

	var (
		total int32
		attrs warehouses.Attributes
		null  sql.NullInt32
	)

	row := s.Conn.QueryRow(context.Background(), query, warehouseID)

	err := row.Scan(&attrs.Name, &attrs.Availability, &null)

	if !null.Valid {
		return warehouses.Attributes{}, 0, fmt.Errorf("%s", errs.WarehouseError)
	}

	if err != nil {
		return warehouses.Attributes{}, 0, err
	}

	total = null.Int32

	return attrs, total, nil
}
