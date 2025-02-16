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

func TestLogin(t *testing.T) {
	a := app.NewApp()
	err := a.Init()
	if err != nil {
		t.Fatalf("Can't init app: %v", err)
	}

	// Мокируем запрос
	loginData := ds.AuthRequest{
		Username: "vvjke314",
		Password: "hello",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Error marshalling login data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Запуск теста
	w := httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp["token"])
}
