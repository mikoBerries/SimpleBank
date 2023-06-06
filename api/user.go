package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

type newUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// createAccount create new account with 0 balance
func (server *server) createUser(ctx *gin.Context) {
	var req newUserRequest
	//Unmarshal request with tag validation
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	saltedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: saltedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok { //if it is pq error
			switch pqErr.Code.Name() {
			case "unique_violation":
				//violation constrain key error in pq
				ctx.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"acc_token"`
	AccesTokenExpiresAt   time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *server) userLogin(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//check naked password and hashed in db
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	accessToken, accessTokenPayload, err := server.token.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Generate Token for session
	refreshToken, refeshTokenPayload, err := server.token.CreateToken(user.Username, server.config.RefeshTokenDuration)
	//Token ID and ExpiredAt is in inside token payload
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//convert to uuid
	uuid, err := uuid.FromString(refeshTokenPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateSessionParams{
		ID:           uuid,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refeshTokenPayload.ExpiresAt.Time,
	}
	resultSession, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := loginUserResponse{
		SessionId:             resultSession.ID,
		AccessToken:           accessToken,
		AccesTokenExpiresAt:   accessTokenPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refeshTokenPayload.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
