package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/MikoBerries/SimpleBank/db/mock"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// eqCreateUserParam costum struct of gomock.Mathcer
type eqCreateUserParam struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParam) Matches(x interface{}) bool {
	//convert/assert x interface{} to expected struct returning false if failed
	actualArg, ok := x.(db.CreateUserParams)

	if !ok { //miss-match
		return false
	}
	err := util.CheckPassword(e.password, actualArg.HashedPassword)
	if err != nil { //check hash password
		return false
	}

	//complete e.arg hashed password
	e.arg.HashedPassword = actualArg.HashedPassword
	//compare it with expected and actual arg
	return reflect.DeepEqual(e.arg, actualArg)
}

func (e eqCreateUserParam) String() string {
	//String of error that mathcer produce if Matches returning false
	failString := fmt.Sprintf("failed to matching create user param.\n want:%v  passwowrd :%v", e.arg, e.password)
	return failString
}

// CreateEqCreateUserParam eqCreateUserParam builder (gomock.Matcher)
func CreateEqCreateUserParam(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParam{arg: arg, password: password}

}

func TestCreateUser(t *testing.T) {
	//Expected user to test
	expectUser, password := randomUser(t)

	ctrl := gomock.NewController(t)

	//make struct to support multi stub testing
	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username":  expectUser.Username,
				"password":  password,
				"full_name": expectUser.FullName,
				"email":     expectUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: expectUser.Username,
					FullName: expectUser.FullName,
					Email:    expectUser.Email,
				}

				store.
					EXPECT().
					CreateUser(gomock.Any(), CreateEqCreateUserParam(arg, password)).
					Times(1).
					Return(expectUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				// matchAccount(t, recorder.Body, expectAccount)
			},
		},
	}
	//loop every test case to run separtly
	for _, tc := range testCase {

		// use t.run to make separate go routine func
		t.Run(tc.name, func(t *testing.T) { //testing name and what testing will execute?
			store := mockdb.NewMockStore(ctrl)

			// access buildstubs function and run it using store as argument
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//make http.NewRequest params
			url := "/createUser"
			byteData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			reqeust, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(byteData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, reqeust)
			tc.checkResponse(t, recorder)
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

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusOK, recorder.Code)
// 			matchAccount(t, recorder.Body, expectAccount)
// 		},
// 	},
// 	{
// 		name:      "notFound",
// 		accountId: expectAccount.ID,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.
// 				EXPECT().
// 				GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
// 				Times(1).
// 				Return(db.Account{}, sql.ErrNoRows)

// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusNotFound, recorder.Code)
// 			// matchAccount(t, recorder.Body, expectAccount)
// 		},
// 	},
// 	{
// 		name:      "internalServerError",
// 		accountId: expectAccount.ID,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.
// 				EXPECT().
// 				GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
// 				Times(1).
// 				Return(db.Account{}, sql.ErrConnDone)

// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			// matchAccount(t, recorder.Body, expectAccount)
// 		},
// 	},
// 	{
// 		name:      "badReqeust",
// 		accountId: -1,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.
// 				EXPECT().
// 				GetAccount(gomock.Any(), gomock.Any()).
// 				Times(0)
// 			// Return(db.Account{}, sql.ErrNoRows)

// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			// matchAccount(t, recorder.Body, expectAccount)
// 		},
// 	},
// }
// //loop each test case and run each test case with data from anymous struct
// for i := range testCase {
// 	tc := testCase[i]
// 	// use t.run to make separate go routine func
// 	t.Run(tc.name, func(t *testing.T) { //testing name and what testing will execute?
// 		store := mockdb.NewMockStore(ctrl)
// 		// access buildstubs function and run it using store as argument
// 		tc.buildStubs(store)

// 		server := NewTestServer(t,store)
// 		recorder := httptest.NewRecorder()

// 		url := fmt.Sprintf("/account/%d", tc.accountId)

// 		reqeust, err := http.NewRequest(http.MethodGet, url, nil)
// 		require.NoError(t, err)
// 		server.router.ServeHTTP(recorder, reqeust)
// 		tc.checkResponse(t, recorder)
// 	})
// }
// //single Test case
// //make mock reflect of Store Stuct
// store := mockdb.NewMockStore(ctrl)
// //one single stub
// store.
// 	EXPECT().
// 	GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
// 	Times(1).
// 	Return(expectAccount, nil)

// //create server to send test request to API
// //get gin router using store DB conn
// server := NewTestServer(t,store)
// recorder := httptest.NewRecorder()

// //url API to test request
// url := fmt.Sprintf("/account/%d", expectAccount.ID)
// //make get request
// reqeust, err := http.NewRequest(http.MethodGet, url, nil)
// require.NoError(t, err)

// server.router.ServeHTTP(recorder, reqeust)
// //all callback recorded in recorder
// //check all callback testing
// require.Equal(t, http.StatusOK, recorder.Code)
// matchAccount(t, recorder.Body, expectAccount)
