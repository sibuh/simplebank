package api

import (
	db "assignment_01/simplebank/db/sqlc"
	"assignment_01/simplebank/util"
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
type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	ChangedPasswordAt time.Time `json:"changed_password_at"`
	CreatedAt         time.Time `json:"created_at"`
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
	resp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		ChangedPasswordAt: user.ChangedPasswordAt,
	}
	ctx.JSON(http.StatusOK, resp)
}
