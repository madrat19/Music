package main

import (
	"log"
	"music/internal/database"
	"music/internal/server"
	"music/mock"
	"os"
)

func main() {
	// Создаем файл для логов и настраиваем вывод
	file, err := os.OpenFile("../app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)

	// Создаем таблицу в бд
	err = database.InitTable()
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	// Запучкаем мок-сервер music_info
	go mock.RunServer()
	// Запускаем сервер приложения
	server.RunServer()
}
