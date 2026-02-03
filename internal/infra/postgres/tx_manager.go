package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dubininme/xm-assessment/internal/domain/company"
)

type ctxKey string

const txKey ctxKey = "tx"

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

var _ company.TxManager = (*TxManager)(nil)

type TxManager struct {
	db *Db
}

func NewTxManager(db *Db) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)
	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func ExtractExecutor(ctx context.Context, db *Db) Executor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return db.DB
}
