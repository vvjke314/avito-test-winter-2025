package app

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("your_secret_key")

func (a App) Login(c *gin.Context) {
	// Получаем username из тела запроса
	var loginRequest ds.AuthRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		resp := ds.NewErrorResponse("Invalid username format")
		resp.Response(c, 400)
		return
	}

	// Проверяем, существует ли пользователь
	employee, err := a.Repo.FindEmployeeByUsername(loginRequest.Username)
	if err != nil {
		resp := ds.NewErrorResponse("Error occurred while searching for user")
		resp.Response(c, 500)
		return
	}

	// Если пользователь не найден, создаем нового
	if employee == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			resp := ds.NewErrorResponse("Error occurred while hashing pass")
			resp.Response(c, 500)
			return
		}

		// Создание нового пользователя
		employee = &ds.Employee{
			Id:           uuid.New(),
			Username:     loginRequest.Username,
			PasswordHash: string(hashedPassword),
			Coins:        1000,
		}
		if err := a.Repo.CreateEmployee(employee); err != nil {
			resp := ds.NewErrorResponse("Error occurred while creating user")
			resp.Response(c, 500)
			return
		}
	}

	// Генерация JWT токена
	tokenString, err := generateToken(employee.Username)
	if err != nil {
		resp := ds.NewErrorResponse("Can't generate JWT-token")
		resp.Response(c, 500)
		return
	}

	// Устанавливаем JWT токен в cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	// Возвращаем токен в ответе
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
