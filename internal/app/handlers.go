package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

func (a *App) BuyItem(c *gin.Context) {
	// Получение имени товара из URL
	itemName := c.Param("item")

	// Получение имени пользователя из контекста
	username, exists := c.Get("username")
	if !exists {
		resp := ds.NewErrorResponse("Unauthorized")
		resp.Response(c, 401)
		return
	}
	usernameStr, _ := username.(string)

	// Поиск пользователя
	user, err := a.Repo.FindEmployeeByUsername(usernameStr)
	if err != nil {
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}

	// Поиск товара
	item, err := a.Repo.FindItemByName(itemName)
	if err != nil {
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}
	if item == nil {
		resp := ds.NewErrorResponse("Bad request")
		resp.Response(c, 400)
		return
	}

	// Проверка баланса
	if user.Coins < item.Price {
		resp := ds.NewErrorResponse("Insufficient funds")
		resp.Response(c, 400)
		return
	}

	// Начало транзакции
	tx, err := a.Repo.DB.Begin()
	if err != nil {
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}

	// Списание монет
	if err := a.Repo.DecreaseUserCoins(tx, user.Id.String(), item.Price); err != nil {
		tx.Rollback()
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}

	// Запись покупки в историю
	if err := a.Repo.RecordPurchase(tx, user.Id.String(), item.Id.String(), item.Price); err != nil {
		tx.Rollback()
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}

	// Фиксация транзакции
	if err := tx.Commit(); err != nil {
		resp := ds.NewErrorResponse("Internal server error")
		resp.Response(c, 500)
		return
	}

	// Ответ клиенту
	c.JSON(http.StatusOK, gin.H{})
}

// SendCoin отправка монет пользователю
func (a *App) SendCoin(c *gin.Context) {
	var request ds.SendCoinRecord

	// Чтение тела запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Получение имени пользователя из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	usernameStr, _ := username.(string)

	// Поиск пользователя-отправителя
	sender, err := a.Repo.FindEmployeeByUsername(usernameStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if sender == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}

	// Поиск получателя
	receiver, err := a.Repo.FindEmployeeByUsername(request.ToUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if receiver == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
		return
	}

	// Проверка наличия достаточных монет у отправителя
	if sender.Coins < request.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Начало транзакции
	tx, err := a.Repo.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Списание монет у отправителя
	if err := a.Repo.DecreaseUserCoins(tx, sender.Id.String(), request.Amount); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrease sender's coins"})
		return
	}

	// Добавление монет получателю
	if err := a.Repo.IncreaseUserCoins(tx, receiver.Id.String(), request.Amount); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increase receiver's coins"})
		return
	}

	// Запись транзакции
	if err := a.Repo.RecordCoinTransfer(tx, sender.Id.String(), receiver.Id.String(), request.Amount); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record coin transfer"})
		return
	}

	// Фиксация транзакции
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Ответ клиенту
	c.JSON(http.StatusOK, gin.H{"message": "Coins successfully sent"})
}

// GetUserInfo получение информации о пользователе (монеты, инвентарь, история переводов)
func (a *App) GetUserInfo(c *gin.Context) {
	// Получение имени пользователя из контекста (из токена)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	usernameStr, _ := username.(string)

	// Получение информации о пользователе
	info, err := a.Repo.GetClientInfo(usernameStr)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Возврат информации в формате JSON
	c.JSON(http.StatusOK, info)
}
