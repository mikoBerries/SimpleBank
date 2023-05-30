package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store Embedding Queries and sq.DB
//to make a struct that can call queries and db to do transaction
type Store struct {
	queires *Queries
	db      *sql.DB
}

//NewStore make and return new  Store struct
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		queires: New(db),
	}
}

//execTx execute transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//Begin Transaction using context and isolation rule
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error when begin transcation %w", err)
	}
	// defer tx.Rollback()
	q := New(tx)
	err = fn(q)
	if err != nil {
		// tx.Rollback()
		if tx.Rollback().Error() != "" {
			return fmt.Errorf("eror when executing %w error when rollback transaction %w", err, tx.Rollback())
		}
		return err
	}
	//if all pass commit changes
	return tx.Commit()
}

//TransferTxParams contains of input param needed to do TransferTX
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult contains of result of TransferTx
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account_id"`
	ToAccount   Account  `json:"to_account_id"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx performs transfer moeny from one account to other account
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		return err
	})

	return result, err
}
