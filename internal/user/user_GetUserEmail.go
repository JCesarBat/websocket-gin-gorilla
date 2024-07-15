package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
	db "webSocketGorrilaMuxGrpc/db/sqlc"
)

const (
	secretKey = "secret"
)

type LoginRequest struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	User        db.User `json:"user"`
	AccessToken string  `json:"access_token" `
}

type MyJwtClaims struct {
	ID       string `json:"id"`
	Username string `json:"username" `
	jwt.RegisteredClaims
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("the request is incorrect %s", err.Error()).Error()})
		return
	}

	user, err := h.store.GetUserEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("not found the user  %s", err.Error()).Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("internal server error  %s", err.Error()).Error()})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJwtClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Errorf("error in token %s", err.Error()).Error()})
		return
	}
	ctx.SetCookie("jwt", ss, 3600, "/", "localhost", false, true)
	response := LoginResponse{
		User:        user,
		AccessToken: ss,
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout susses"})
}
