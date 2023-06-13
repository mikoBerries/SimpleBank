package gapi

import (
	"context"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/val"
	"github.com/MikoBerries/SimpleBank/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser to create new username and serving for gRPC server
func (server *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	//Check all violations happend on this request
	violations := validateCreateUserReqeust(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	saltedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: saltedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		//build stub for this create user process
		AfterCreate: func(user db.User) error {
			//prepare payload for redis
			payload := &worker.PayloadSendVerifyEmail{Username: user.Username}
			// set list of this task option
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(2 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			// distribute task with payload and opts
			err = server.taskDisributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)

			// log.Info().Str("msg", "Task sended").Send()
			return err
		},
	}

	txUserResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		//mapped erorr in postgres
		if db.ErrorCode(err) == db.UniqueViolation {
			//returning http/2 codes violations instead pq err code
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txUserResult.User),
	}
	return rsp, nil
}

// ValidateCreateUserReqeust return all violation (code: BadReqeust) happended at create function reqeust
func validateCreateUserReqeust(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//check everything nedeed
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return
}
