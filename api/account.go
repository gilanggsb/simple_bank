package api

import (
	"database/sql"
	db "gilanggsb/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, http.StatusBadRequest))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, http.StatusBadRequest))
		return
	}

	ctx.JSON(http.StatusOK, baseResponse(&BaseResponse[db.Accounts]{
		Status:  http.StatusOK,
		Data:    account,
		Message: "Success create account",
	}))
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, http.StatusBadRequest))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, http.StatusBadRequest))
		return
	}

	ctx.JSON(http.StatusOK, baseResponse(&BaseResponse[db.Accounts]{
		Status:  http.StatusOK,
		Data:    account,
		Message: "Success get account ",
	}))
}

type listAccountRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, http.StatusBadRequest))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  int32(req.PageSize),
		Offset: int32(req.PageID-1) + int32(req.PageSize),
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, http.StatusBadRequest))
		return
	}

	ctx.JSON(http.StatusOK, baseResponse(&BaseResponse[[]db.Accounts]{
		Status:  http.StatusOK,
		Data:    accounts,
		Message: "Success get list accounts",
	}))
}
