package api

import (
	"database/sql"
	"fmt"
	db "gilanggsb/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required",min=1`
	ToAccountID   int64  `json:"to_account_id" binding:"required",min=1`
	Amount        int64  `json:"amount" binding:"required",gt=0`
	Currency      string `json:"currency" binding:"required",currency`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, http.StatusBadRequest))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		// internal error should respond with 500 both in HTTP code and payload status
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, baseResponse(&BaseResponse[db.TransferTxResult]{
		Status:  http.StatusOK,
		Data:    transfer,
		Message: "Success create transfer",
	}))
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found
			ctx.JSON(http.StatusNotFound, errorResponse(err, http.StatusNotFound))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err, http.StatusInternalServerError))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency missmatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err, http.StatusBadRequest))
		return false
	}

	return true

}
