package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/token"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	// Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

//createAccount create new account with 0 balance
func (server *server) createAccount(ctx *gin.Context) {
	//Unmarshal request with tag validation
	var req CreateAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	//ctx assert to token.payload
	authPaylaod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPaylaod.UserName,
		Balance:  0,
		Currency: req.Currency,
	}

	acc, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, acc)
}

type AccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,numeric,min=1" `
}

// getAccountByID getting 1 account data by ID
func (server *server) getAccountByID(ctx *gin.Context) {
	var req AccountByIDRequest
	//Unmarshal request with tag validation
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	acc, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		//error now data retirfe
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	//ctx assert to token.payload
	authPaylaod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if acc.Owner != authPaylaod.UserName { //requester are not owner
		err = errors.New("account doesn`t belong to auth owner")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type ListAccountRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	ItemPages  int32 `form:"item_pages" binding:"required,min=5,max=10"`
}

// getAccountByID getting 1 account data by ID
func (server *server) getListAccount(ctx *gin.Context) {
	var req ListAccountRequest

	//Unmarshal request with tag validation
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	arg := db.ListAccountsParams{
		Limit:  req.ItemPages,
		Offset: (req.PageNumber - 1) * req.ItemPages,
	}

	acc, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if len(acc) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"err": "Empty"})
		return

	}
	ctx.JSON(http.StatusOK, acc)
}
