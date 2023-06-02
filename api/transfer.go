package api

import (
	"database/sql"
	"net/http"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,IsCurrency"`
}

//TransferTx create transfer data
func (server *server) createTransfer(ctx *gin.Context) {
	var req TransferRequest
	//binding checking tags
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		//error no data
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		//other error
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
