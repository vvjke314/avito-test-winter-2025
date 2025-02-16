package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func NewApp() *App {
	return &App{}
}

func (a App) Init() {
	
}
