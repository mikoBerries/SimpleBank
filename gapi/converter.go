package gapi

import (
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertUser convert db.User to pb.User
func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreateAt:          timestamppb.New(user.CreatedAt),
	}
}
