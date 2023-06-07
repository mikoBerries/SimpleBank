package gapi

import (
	"context"
	"database/sql"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LoginUser serve GRPC func for login user rpc
func (server *server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")

		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	//check naked password and hashed in db
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "incorrect password")
	}
	accessToken, accessTokenPayload, err := server.token.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	//Generate Token for session
	refreshToken, refeshTokenPayload, err := server.token.CreateToken(user.Username, server.config.RefeshTokenDuration)
	//Token ID and ExpiredAt is in inside token payload
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	//convert to uuid
	uuid, err := uuid.FromString(refeshTokenPayload.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert uuid")
	}
	mtaData := server.extractMetadata(ctx)
	arg := db.CreateSessionParams{
		ID:           uuid,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtaData.UserAgent,
		ClientIp:     mtaData.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refeshTokenPayload.ExpiresAt.Time,
	}
	resultSession, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create seesion")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             resultSession.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refeshTokenPayload.ExpiresAt.Time),
	}

	return rsp, nil
}
