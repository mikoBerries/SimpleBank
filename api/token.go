package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken         string    `json:"access_token"`
	AccesTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// renewAccessToken require RefreshToken produce a new AccessToken with new expiresAt
// used by client after server responing AccessToken are expired
// or client side token are (time.now > AccesTokenExpiresAt)
func (server *server) renewAccessTokenRequest(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//verify RefreshToken are decrpit-able ? / are payload.Valid() ?
	refeshPayload, err := server.token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Check incoming Refreshtoken id(UUID) to Database session
	session, err := server.store.GetSession(ctx, refeshPayload.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Check if this session with this RenewToken.Payload.ID are blocked ? (Force to login again to produce new session)
	if session.IsBlocked {
		err = errors.New("session are Blocked")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err = errors.New("missmatch refreshtoken")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//check and compare refreshtoken payload with session data
	//paylaod username
	if session.Username != refeshPayload.UserName {
		err = errors.New("incorrect user session")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//Generate Token for session
	accessToken, accessTokenPayload, err := server.token.CreateToken(session.Username, server.config.AccessTokenDuration)
	//Token ID and ExpiredAt is in inside token payload
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:         accessToken,
		AccesTokenExpiresAt: accessTokenPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, rsp)
}
