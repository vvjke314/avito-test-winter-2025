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
	username := c.PostForm("username")

	// Поиск пользователя в БД

	tokenString, err := generateToken(username)
	if err != nil {
		resp := ds.NewErrorResponse()
		resp.Logout(c, 500)
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
