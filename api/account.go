package api

import (
	"database/sql"
	"net/http"

	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// createAccountRequest is the request body for createAccount API.
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount creates a new account with a given owner and currency from the request.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorReponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// getAccountRequest is the request body for getAccount API.
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount gets the account info from the request.
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
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

// listAccountsRequest is the request body for listAccount API.
type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts lists all accounts from the request.
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
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
