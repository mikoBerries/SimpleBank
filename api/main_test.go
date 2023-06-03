package api

import (
	"os"
	"testing"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

//newTestServer make new Server struct for test enviroment
func newTestServer(t *testing.T, store db.Store) *server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}
