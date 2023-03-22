package api

import (
	"errors"
	"exercise/simplebank/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authType                = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)
		if authorization == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("no authorization"))
			return
		}
		fields := strings.Fields(authorization)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("invalid authorization format"))
			return
		}
		authPrefix := strings.ToLower(fields[0])
		if authPrefix != authType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unsupported authorization type"))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("incorrect access token"))
			return
		}
		err = payload.Valid()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
