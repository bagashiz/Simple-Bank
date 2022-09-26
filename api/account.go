package api

import (
	"database/sql"
	"net/http"

	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateAccountRequest is the request body for createAccount API.
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// CreateAccount creates a new account with a given owner and currency from the request.
func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// GetAccountRequest is the request body for getAccount API.
type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount gets the account info from the request.
func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorReponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// ListAccountRequest is the request body for listAccount API.
type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts lists all accounts from the request.
func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
