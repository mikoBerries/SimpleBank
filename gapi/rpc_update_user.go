package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser to create new username and serving for gRPC server
func (server *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	//checking header token auth (Paseto v2)
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	//Check all violations happend on this request
	violations := validateUpdateUserReqeust(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.UserName != req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{String: req.GetFullName(), Valid: req.GetFullName() != ""},
		Email:    sql.NullString{String: req.GetEmail(), Valid: req.GetEmail() != ""},
	}

	//check password request separate
	if req.GetPassword() != "" {
		saltedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
		}
		arg.HashedPassword = sql.NullString{String: saltedPassword, Valid: true}
		arg.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, status.Errorf(codes.NotFound, "user not found")
		case sql.ErrConnDone:
			return nil, status.Errorf(codes.Internal, "failed connection done")
		case sql.ErrTxDone:
			return nil, status.Errorf(codes.Internal, "failed trx done")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get user")
		}
	}
	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

// ValidateCreateUserReqeust return all violation (code: BadReqeust) happended at create function reqeust
func validateUpdateUserReqeust(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//check everything nedeed
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	updateData := false
	if req.GetPassword() != "" {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
		updateData = true
	}

	if req.GetFullName() != "" {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
		updateData = true
	}

	if req.GetEmail() != "" {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
		updateData = true
	}
	//nothings to update
	if !updateData {
		violations = append(violations, fieldViolation("field", fmt.Errorf("missing field to update")))
	}
	return
}
