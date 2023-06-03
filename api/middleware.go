package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MikoBerries/SimpleBank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authoriztion_payload"
)

//authMiddleWare middleware to auth token in header
func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get from header of request
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 { //no auth header get from request header
			err := errors.New("authorization Header missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 { //header must contain 2 field auth type and token("Bearer" and "paseto token")
			err := errors.New("invalid authorization Header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer { //auth type not bearer
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return

		}
		//check paseto token with public key
		accesToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accesToken)
		if err != nil { //token failed to verify
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		//set payload to context so payload cam used by other function.handler by context
		ctx.Set(authorizationPayloadKey, payload)

		//done for authMiddleWare and continued to actual gin.router path
		ctx.Next()
	}
}
