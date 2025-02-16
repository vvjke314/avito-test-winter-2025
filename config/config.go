package config

import (
	"fmt"
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
func LoadConfig() (DBConfig, error) {
	host := os.Getenv("DB_HOST")

	port := os.Getenv("DB_PORT")

	user := os.Getenv("DB_USER")

	password := os.Getenv("DB_PASSWORD")

	dbName := os.Getenv("DB_NAME")

	// Валидация переменных окружения

	if host == "" {

		return DBConfig{}, fmt.Errorf("DB_HOST is required")

	}

	if port == "" {

		return DBConfig{}, fmt.Errorf("DB_PORT is required")

	}

	if user == "" {

		return DBConfig{}, fmt.Errorf("DB_USER is required")

	}

	if password == "" {

		return DBConfig{}, fmt.Errorf("DB_PASSWORD is required")

	}

	if dbName == "" {

		return DBConfig{}, fmt.Errorf("DB_NAME is required")

	}

	return DBConfig{

		Host: host,

		Port: port,

		User: user,

		Password: password,

		DBName: dbName,
	}, nil
}
