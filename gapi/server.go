package gapi

import (
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/token"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/worker"
)

// Server serve GRPC request for apps sevices.
type server struct {
	pb.UnimplementedSimplebankServer
	store          db.Store
	token          token.Maker
	config         util.Config
	taskDisributor worker.TaskDistributor
}

// NewServer Create new GRPC Server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*server, error) {
	//crete tokeMaker and sign it to server
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &server{
		store:          store,
		token:          tokenMaker,
		config:         config,
		taskDisributor: taskDistributor,
	}
	return server, nil
}
