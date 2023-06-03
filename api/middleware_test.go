package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MikoBerries/SimpleBank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorizationHeader(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//Build auth header value (authType token) and set request header
	authHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authHeader)
}

func TestAuthMiddleWare(t *testing.T) {
	testCase := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Token Expired",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization Header",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization Header type",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupport Authorization Header type",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, "something", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// matchAccount(t, recorder.Body, expectAccount)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		// use t.run to make separate go routine func
		t.Run(tc.name, func(t *testing.T) { //testing name and what testing will execute?
			//setup new server just only to check middle ware func
			server := newTestServer(t, nil)
			authUrl := "/auth"
			server.router.GET(authUrl, authMiddleWare(server.token),
				func(ctx *gin.Context) { //after ctx.next are called
					ctx.JSON(http.StatusOK, gin.H{})
				})

			// server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//create request and set header for each test case
			reqeust, err := http.NewRequest(http.MethodGet, authUrl, nil)
			require.NoError(t, err)
			tc.setupAuth(t, reqeust, server.token)

			//execute request and check response
			server.router.ServeHTTP(recorder, reqeust)
			tc.checkResponse(t, recorder)
		})
	}

}
