package gapi

import (
	"testing"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/worker"
	"github.com/stretchr/testify/require"
)

// newTestServer make new Server struct for test enviroment (gRPC server)
func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	//
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}

// func TestMain(m *testing.M) {

// 	os.Exit(m.Run())
// }
