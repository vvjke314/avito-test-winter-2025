package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

var secretKey = []byte("vvjke314")

func Login(c *gin.Context) {
	// Начало запроса
	var username string
	// Начать обрабатывать запрос

	tokenString, err := generateToken(username)
	if err != nil {
		resp := ds.NewErrorResponse("Can't generate JWT-token")
		resp.Response(c, 500)
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
