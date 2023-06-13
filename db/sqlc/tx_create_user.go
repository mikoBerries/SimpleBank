package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

//CreataUserTx create user with transaction to rollback
func (sqlStore *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := sqlStore.execTx(ctx, func(q *Queries) error {
		var err error
		//create user first
		result.User, err = sqlStore.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		//do injected function after create user success
		err = arg.AfterCreate(result.User)
		//will returning error if function not satisfied
		return err
	})

	return result, err
}
