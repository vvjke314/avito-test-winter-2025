package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vvjke314/avito-test-winter-2025/config"
	"github.com/vvjke314/avito-test-winter-2025/internal/repository"
)

type App struct {
	Router *gin.Engine
	Repo   repository.Repo
}

func NewApp() *App {
	return &App{}
}

// Init инициализирует приложения, подключаясь к бд и проводя миграции по инициализации приложения
func (a *App) Init() error {
	dbConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("[config.LoadConfig] can't load env variables, %w", err)
	}
	connStr := "postgres://" + dbConfig.User + ":" + dbConfig.Password + "@" + dbConfig.Host + ":" + dbConfig.Port + "/" + dbConfig.DBName

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("[sql.Open] can't connect to database, %w", err)
	}
	a.Repo.DB = conn

	if err := goose.Up(a.Repo.DB, "migrations"); err != nil {
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
