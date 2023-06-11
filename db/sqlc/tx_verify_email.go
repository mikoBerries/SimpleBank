package db

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	EmailId     int64
	SecreteCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// CreataUserTx create user with transaction to rollback
func (sqlStore *SqlStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	err := sqlStore.execTx(ctx, func(q *Queries) error {
		var err error
		//Update verify email data
		arg := UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecreteCode,
		}
		//get result from updating VerifyEmail
		result.VerifyEmail, err = sqlStore.UpdateVerifyEmail(ctx, arg)
		if err != nil {
			return err
		}
		//just update is emailverified field the rest of parma is empty not updated
		argUpdateUser := UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		}
		//get result from updating user IsEmailVerified to true
		result.User, err = sqlStore.UpdateUser(ctx, argUpdateUser)
		if err != nil {
			return err
		}
		return err
	})

	return result, err
}
