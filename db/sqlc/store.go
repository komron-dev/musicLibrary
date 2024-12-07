package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Store provides all functions to execute db queries and transaction
type Store struct {
	db *sql.DB
	*Queries
	logger *logrus.Logger
}

// NewStore creates a new store
func NewStore(db *sql.DB, logger *logrus.Logger) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
		logger:  logger,
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			store.logger.Error("Failed to rollback:", rbErr)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
