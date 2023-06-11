package gapi

import (
	"context"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyEmail to create new username and serving for gRPC server
func (server *server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	//Check all violations happend on this request
	violations := validateVerifyEmailReqeust(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	arg := db.VerifyEmailTxParams{
		EmailId:     req.GetEmailId(),
		SecreteCode: req.GetSecretCode(),
	}
	verifyResult, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}
	rsp := &pb.VerifyEmailResponse{
		IsVerified: verifyResult.User.IsEmailVerified,
	}
	return rsp, nil
}

// ValidateVerifyEmailReqeust return all violation (code: BadReqeust) happended at create function reqeust
func validateVerifyEmailReqeust(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//check everything nedeed
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return
}
