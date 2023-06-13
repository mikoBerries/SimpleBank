package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/token"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// newTestServer make new Server struct for test enviroment (gRPC server)
func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}

// newContextWithBearerToken code refactoring used for embeding Berare token auth in incoming context/request for test purpose
func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	ctx := context.Background()
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	//metada struct just map[string]string
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(ctx, md)
}
