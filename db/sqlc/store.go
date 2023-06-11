package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store interface for mock testing need
type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreataUserTx(ctx context.Context, arg CreataUserTxParams) (CreataUserTxResult, error)
	Querier
}

// Store Embedding Queries and sq.DB
// to make a struct that can call queries and db to do transaction
type SqlStore struct {
	db *sql.DB
	*Queries
}

// NewStore make and return new  Store struct
func NewStore(db *sql.DB) *SqlStore {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx Wrapped sqlStore with sql-transcation to rollback when returning err / commit when it's nil
func (sqlStore *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	//Begin Transaction using context and isolation rule
	tx, err := sqlStore.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error when begin transcation %w", err)
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("eror when executing %v error when rollback transaction %v", err, rbErr)
		}
		return err
	}
	//if all pass commit changes
	return tx.Commit()
}
