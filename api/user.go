package api

import (
	db "assignment_01/simplebank/db/sqlc"
	"assignment_01/simplebank/util"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username"  binding:"required,alphanum"`
	Fullname string `json:"full_name"  binding:"required"`
	Password string `json:"password"  binding:"required,min=6"`
	Email    string `json:"email"  binding:"required,email"`
}
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	ChangedPasswordAt time.Time `json:"changed_password_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		ChangedPasswordAt: user.ChangedPasswordAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	HashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: HashedPassword,
		FullName:     req.Fullname,
		Email:        req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}

		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	resp := NewUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Username string `json:"username"  binding:"required,alphanum"`

	Password string `json:"password"  binding:"required,min=6"`
}
type userLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	fmt.Println(req)
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	match := util.CheckPasswordHash(req.Password, user.HashPassword)
	if !match {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	resp := userLoginResponse{
		AccessToken: accessToken,
		User:        NewUserResponse(user),
	}
	fmt.Println(resp)
	ctx.JSON(http.StatusOK, resp)

}
