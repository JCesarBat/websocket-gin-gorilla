package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	db "webSocketGorrilaMuxGrpc/db/sqlc"
)

type Handler struct {
	store  db.Store
	router *gin.Engine
}

func NewHandler(store db.Store) *Handler {

	return &Handler{
		store: store,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	User db.User `json:"user"`
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	param := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := h.store.CreateUser(context.Background(), param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := CreateUserResponse{
		user,
	}
	ctx.JSON(http.StatusOK, response)
}
