package transactionManager

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) TxBegin(ctx context.Context) (pgx.Tx, error) {
	tx, err := tm.db.Begin(ctx)
	return tx, err
}

func (tm *TransactionManager) TxRollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	return err
}

func (tm *TransactionManager) TxCommit(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	return err
}
