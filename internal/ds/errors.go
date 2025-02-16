package ds

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse модель данных для отправки ответов с ошибками клиенту
type ErrorResponse struct {
	Errors string `json:"errors"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Errors: message,
	}
}

// Response фукнкция для отправки ответа ошибки клиенту
func (er *ErrorResponse) Response(c *gin.Context, code int) {
	c.JSON(http.StatusBadRequest, er)
}
