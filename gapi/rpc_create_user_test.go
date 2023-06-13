package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/MikoBerries/SimpleBank/db/mock"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/MikoBerries/SimpleBank/worker"
	mockworker "github.com/MikoBerries/SimpleBank/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// eqCreateUserTxParam costum struct of gomock.Mathcer
type eqCreateUserTxParam struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParam) Matches(x interface{}) bool {
	//convert/assert x interface{} to expected struct returning false if failed
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok { //miss-match
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil { //check hash password
		return false
	}

	//complete e.arg hashed password
	expected.arg.HashedPassword = actualArg.HashedPassword
	//compare it with expected and actual arg param (since embeded func always same)
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	//after create user are matches then executed embeded after func here
	if err = actualArg.AfterCreate(expected.user); err != nil {
		return false
	}
	return true

}

func (expected eqCreateUserTxParam) String() string {
	//String of error that mathcer produce if Matches returning false
	failString := fmt.Sprintf("failed to matching create user param.\n want:%v  passwowrd :%v", expected.arg, expected.password)
	return failString
}

// CreateEqCreateUserTxParam eqCreateUserParam builder (gomock.Matcher)
func CreateEqCreateUserTxParam(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParam{arg: arg, password: password, user: user}

}

func TestCreateUser(t *testing.T) {
	//Expected user to test
	expectedUser, password := randomUser(t)

	//make struct to support multi stub testing
	testCase := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
		checkResponse func(t *testing.T, resp *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: expectedUser.Username,
				Password: password,
				FullName: expectedUser.FullName,
				Email:    expectedUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				//rule expected called function in store interface
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: expectedUser.Username,
						FullName: expectedUser.FullName,
						Email:    expectedUser.Email,
					},
				}
				store.
					EXPECT().
					CreateUserTx(gomock.Any(), CreateEqCreateUserTxParam(arg, password, expectedUser)).
					Times(1).
					Return(db.CreateUserTxResult{User: expectedUser}, nil)

				//rule expected called function in taskDistributor interface
				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: expectedUser.Username,
				}
				taskDistributor.
					EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)

			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)

				createdUser := resp.GetUser()
				//compare expected with actual created user
				require.Equal(t, expectedUser.Username, createdUser.Username)
				require.Equal(t, expectedUser.FullName, createdUser.FullName)
				require.Equal(t, expectedUser.Email, createdUser.Email)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: expectedUser.Username,
				Password: password,
				FullName: expectedUser.FullName,
				Email:    expectedUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				//rule expected called function in store interface

				store.
					EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				//rule expected called function in taskDistributor interface
				taskDistributor.
					EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				//extract error massage
				st, ok := status.FromError(err)
				require.True(t, ok)
				//compare code
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "DuplicateUsername",
			req: &pb.CreateUserRequest{
				Username: expectedUser.Username,
				Password: password,
				FullName: expectedUser.FullName,
				Email:    expectedUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, db.ErrUniqueViolation)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.CreateUserRequest{
				Username: expectedUser.Username,
				Password: password,
				FullName: expectedUser.FullName,
				Email:    "this_is/not?and1email.com@",
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
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

			taskDistributorControler := gomock.NewController(t)
			defer taskDistributorControler.Finish()
			taskDistributor := mockworker.NewMockTaskDistributor(taskDistributorControler)

			tc.buildStubs(store, taskDistributor)
			//get server struct (gRPC)
			server := newTestServer(t, store, taskDistributor)
			//feed server function with ctx and request
			resp, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, resp, err)
		})
	}
}

// randomUser create and returning random user and user saltedpassword
func randomUser(t *testing.T) (db.User, string) {
	randomPassowrd := "nakedpass" + util.RandomString(4)
	saltedPassword, err := util.HashPassword(randomPassowrd)
	require.NoError(t, err)
	createdUser := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: saltedPassword,
		FullName:       util.RandomOwner() + util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return createdUser, randomPassowrd
}
