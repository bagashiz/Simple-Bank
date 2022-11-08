package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/bagashiz/Simple-Bank/token"
)

// transferRequest is the request body for createTransfer API.
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// createTransfer creates a new transfer transaction from the request.
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

  fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency) 
	if !valid {
		return
	}

  _, valid = server.validAccount(ctx, req.ToAccountID, req.Currency) 
	if !valid {
		return
	}

  authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
  if fromAccount.Owner != authPayload.Username {
    err := errors.New("from account does not belong to the authenticated user")
    ctx.JSON(http.StatusUnauthorized, errorReponse(err))
  }
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if the account is really exists and the currency matches.
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorReponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return account, false
	}

	return account, true
}
