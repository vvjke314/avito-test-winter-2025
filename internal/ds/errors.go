package ds

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse модель данных для отправки ответов с ошибками клиенту
type ErrorResponse struct {
	Errors string `json:"errors"`
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

// Logout фукнкция для отправки ответа ошибки клиенту
func (er *ErrorResponse) Logout(c *gin.Context, code int) {
	switch code {
	case 400:
		er.Errors = "Bad request"
	case 401:
		er.Errors = "Unauthorized"
	case 500:
		er.Errors = "Internal service error"
	}
	c.JSON(http.StatusBadRequest, er)
}
