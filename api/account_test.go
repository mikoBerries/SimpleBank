package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MikoBerries/SimpleBank/db/mock"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountByID(t *testing.T) {
	//Expected Test account
	expectAccount := createRandomAccount()

	ctrl := gomock.NewController(t)
	//a mockgen version of 1.5.0+, and are passing a *testing.T into gomock.NewController(t)
	// you no longer need to call ctrl.Finish() explicitly
	// defer ctrl.Finish()

	//make struct to support multi stub testing
	testCase := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountId: expectAccount.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
					Times(1).
					Return(expectAccount, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchAccount(t, recorder.Body, expectAccount)
			},
		},
		{
			name:      "notFound",
			accountId: expectAccount.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// matchAccount(t, recorder.Body, expectAccount)
			},
		},
		{
			name:      "internalServerError",
			accountId: expectAccount.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(expectAccount.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// matchAccount(t, recorder.Body, expectAccount)
			},
		},
		{
			name:      "badReqeust",
			accountId: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				// Return(db.Account{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// matchAccount(t, recorder.Body, expectAccount)
			},
		},
	}
	//loop each test case and run each test case with data from anymous struct
	for i := range testCase {
		tc := testCase[i]
		// use t.run to make separate go routine func
		t.Run(tc.name, func(t *testing.T) { //testing name and what testing will execute?
			store := mockdb.NewMockStore(ctrl)
			// access buildstubs function and run it using store as argument
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.accountId)

			reqeust, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, reqeust)
			tc.checkResponse(t, recorder)
		})
	}
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
}

func createRandomAccount() (account db.Account) {
	account.ID = util.RandomInt(10, 10)
	account.Balance = 0
	account.Currency = util.RandomCurrency()
	return
}

func matchAccount(t *testing.T, body *bytes.Buffer, expectAccount db.Account) {
	//body data are json from request.body
	bodyData, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var actualAccount db.Account
	err = json.Unmarshal(bodyData, &actualAccount)

	require.NoError(t, err)
	require.Equal(t, expectAccount, actualAccount)
}
