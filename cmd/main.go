package main

import (
	"log"
	"music/internal/database"
	"music/internal/server"
	"music/mock"
	"music/tools"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Создаем файл для логов и настраиваем вывод
	file, err := os.OpenFile("../app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	//log.SetOutput(file)

	//Миграции
	config := tools.GetConfig()
	db, err := database.OpenConnection(config)
	if err != nil {
		log.Fatal(1, err)
	}

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(2, err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", migrationDriver)
	if err != nil {
		log.Fatal(3, err)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Запучкаем мок-сервер music_info
	go mock.RunServer()
	// Запускаем сервер приложения
	server.RunServer()
}
