package main

import (
	"log"

	"github.com/vvjke314/avito-test-winter-2025/app"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)

		}
	}()
	
	a := app.NewApp()

	// Инициализация приложения
	if err := a.Init(); err != nil {
		log.Fatalf("Failed to initzialize app: %v", err)
	}

	// Запуск приложения
	a.Run(":8080")
}
