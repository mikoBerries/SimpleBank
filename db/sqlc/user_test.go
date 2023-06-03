package db

import (
	"context"
	"testing"

	"github.com/MikoBerries/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	saltedPassword, err := util.HashPassword("SomeNakedPassword" + util.RandomString(4))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: saltedPassword,
		FullName:       util.RandomOwner() + util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := TestQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)

	// require.NotZero(t, user.ID)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}
