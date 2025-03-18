package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TransactionManager
type TransactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager
func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// WithTransaction
func (m *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
			}
		}
	}()

	err = fn(ctx, tx)
	return err
}

// DB
func (m *TransactionManager) DB() *sqlx.DB {
	return m.db
}
