package config

import (
	"os"
)

// DBConfig содержит параметры подключения к базе данных
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadConfig загружает параметры из переменных окружения
func LoadConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}
