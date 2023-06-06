package token

import (
	"testing"
	"time"

	"github.com/MikoBerries/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	randomSecretKey := util.RandomString(32)
	maker, err := NewJWTMaker(randomSecretKey)
	require.NoError(t, err)

	username := "Mikoberries"
	duration := time.Minute
	role := "user"
	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	//do payload func valid() embed in standart jwt claim

	require.NoError(t, err)
	//check other costum payload/claim value
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.UserName)
	require.Equal(t, role, payload.Role)
}

func TestExpiredJWTMaker(t *testing.T) {
	randomSecretKey := util.RandomString(32)
	maker, err := NewJWTMaker(randomSecretKey)
	require.NoError(t, err)

	username := "Mikoberries"
	duration := -1 * time.Minute
	// role := "user"
	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	_, err = maker.VerifyToken(token)
	require.Error(t, err)
	// require.NotEmpty(t, payload)
	//do payload func valid() embed in standart jwt claim
	// err = payload.Valid()
	// require.Error(t, err)
	//check other costum payload/claim value

}
