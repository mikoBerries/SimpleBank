// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	ListAccount(ctx context.Context) ([]Account, error)
	UpdateAccounts(ctx context.Context, arg UpdateAccountsParams) error
}

var _ Querier = (*Queries)(nil)
