package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/MikoBerries/SimpleBank/db/mock"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/token"
	"github.com/MikoBerries/SimpleBank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser(t *testing.T) {
	//target user to test
	targetUser, _ := randomUser(t)

	newFullname := util.RandomString(6)
	newEmail := util.RandomEmail()

	ivalidemail := "invalid-email"

	//make struct to support multi stub testing
	testCase := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, resp *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: targetUser.Username,
				FullName: &newFullname,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, targetUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				//rule expected called function in store interface
				arg := db.UpdateUserParams{
					Username: targetUser.Username,
					FullName: sql.NullString{
						String: newFullname,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}
				//updated struct of user
				updatedUser := db.User{
					Username:          targetUser.Username,
					HashedPassword:    targetUser.HashedPassword,
					PasswordChangedAt: targetUser.PasswordChangedAt,
					CreatedAt:         targetUser.CreatedAt,
					FullName:          newFullname,
					Email:             newEmail,
				}
				store.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)

			},

			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)

				updatedUser := resp.GetUser()
				//compare expected with actual created user
				require.Equal(t, targetUser.Username, updatedUser.Username)
				require.Equal(t, newFullname, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{
			name: "UsernameNotFound",
			req: &pb.UpdateUserRequest{
				Username: targetUser.Username,
				FullName: &newFullname,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, targetUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.Error(t, err)

				//extract error massage
				st, ok := status.FromError(err)
				require.True(t, ok)
				//compare code
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username: targetUser.Username,
				FullName: &newFullname,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, targetUser.Username, -time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.Error(t, err)

				//extract error massage
				st, ok := status.FromError(err)
				require.True(t, ok)
				//compare code
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoTokenAuthorization",
			req: &pb.UpdateUserRequest{
				Username: targetUser.Username,
				FullName: &newFullname,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.Error(t, err)

				//extract error massage
				st, ok := status.FromError(err)
				require.True(t, ok)
				//compare code
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmailRequest",
			req: &pb.UpdateUserRequest{
				Username: targetUser.Username,
				FullName: &newFullname,
				Email:    &ivalidemail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, targetUser.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.Error(t, err)

				//extract error massage
				st, ok := status.FromError(err)
				require.True(t, ok)
				//compare code
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
	}
	//loop every test case to run separtly
	for _, tc := range testCase {

		// use t.run to make separate go routine func
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)
			//get server struct (gRPC)
			server := newTestServer(t, store, nil)
			//build context with header token auth
			ctx := tc.buildContext(t, server.token)
			resp, err := server.UpdateUser(ctx, tc.req)

			tc.checkResponse(t, resp, err)
		})
	}
}
