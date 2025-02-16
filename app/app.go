package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/vvjke314/avito-test-winter-2025/config"
)

type App struct {
	Router *gin.Engine
	DBConn *sql.DB
}

func NewApp() *App {
	return &App{}
}

// Init инициализирует приложения, подключаясь к бд и проводя миграции по инициализации приложения
func (a *App) Init() error {
	dbConfig := config.LoadConfig()
	connStr := "postgres://" + dbConfig.User + ":" + dbConfig.Password + "@" + dbConfig.Host + ":" + dbConfig.Port + "/" + dbConfig.DBName

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("[sql.Open] can't connect to database, %w", err)
	}
	a.DBConn = conn

	if err := goose.Up(a.DBConn, "migrations"); err != nil {
		return fmt.Errorf("[goose.Up] can't make migrations: %w", err)
	}

	a.Router = gin.Default()
	a.setRouting()
	return nil
}

// setRouting устанавливает возможные ручки для сервиса
func (a *App) setRouting() {

}

// Run запускает приложение
func (a *App) Run(port string) error {

	if err := a.Router.Run(port); err != nil {
		return fmt.Errorf("[gin.Engine.Run] :%w", err)
	}

	return nil	
}
