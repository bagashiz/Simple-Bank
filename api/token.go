package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// renewAccessTokenRequest is the request body for renew access token request.
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorReponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorReponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorReponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("session username does not match")
		ctx.JSON(http.StatusUnauthorized, errorReponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("refresh token does not match")
		ctx.JSON(http.StatusUnauthorized, errorReponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("session is expired")
		ctx.JSON(http.StatusUnauthorized, errorReponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
