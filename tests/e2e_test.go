package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vvjke314/avito-test-winter-2025/internal/app"
	"github.com/vvjke314/avito-test-winter-2025/internal/ds"
)

func TestBuyItem(t *testing.T) {
	// Инициализация приложения
	a := app.NewApp()
	err := a.Init()
	if err != nil {
		t.Fatalf("Can't init app: %v", err)
	}

	loginData := ds.AuthRequest{
		Username: "qwerty",
		Password: "vvs",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Error marshalling login data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	token := resp["token"]
	assert.NotEmpty(t, token)

	// Покупка мерча
	req, _ = http.NewRequest("GET", "/api/buy/cup", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, w.Code)
	var buyResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &buyResp)
	assert.Nil(t, err)
	assert.Equal(t, "", buyResp["message"])
}

func TestSendCoin(t *testing.T) {
	// Инициализация приложения
	a := app.NewApp()
	err := a.Init()
	if err != nil {
		t.Fatalf("Can't init app: %v", err)
	}

	// Авторизация (логиним пользователя-отправителя)
	loginData := ds.AuthRequest{
		Username: "qwerty",
		Password: "vvs",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Error marshalling login data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	token := resp["token"]
	assert.NotEmpty(t, token)

	// Передача монеток
	sendCoinData := ds.SendCoinRecord{
		ToUser: "vvjke314",
		Amount: 10,
	}

	jsonData, err = json.Marshal(sendCoinData)
	if err != nil {
		t.Fatalf("Error marshalling send coin data: %v", err)
	}

	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, w.Code)
	var sendCoinResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &sendCoinResp)
	assert.Nil(t, err)
	assert.Equal(t, "Coins successfully sent", sendCoinResp["message"])
}
