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

	// Получаем настройки приложения
	config := tools.GetConfig()

	// Создаем файл для логов и настраиваем вывод
	file, err := os.OpenFile("../app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	level := config.LogLevel
	infoLog := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(file, "ERROR\t", log.Ldate|log.Ltime)
	fatalLog := log.New(file, "FATAL\t", log.Ldate|log.Ltime)
	tools.InitLogger(level, infoLog, errorLog, fatalLog)

	//Миграци
	db, err := database.OpenConnection(config)
	if err != nil {
		tools.Logger.Fatal("Failed to open db connection: ", err)
	}

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		tools.Logger.Fatal("Failed to get migration driver: ", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", migrationDriver)
	if err != nil {
		tools.Logger.Fatal("Failed to get  migrator: ", err)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		tools.Logger.Fatal("Failed to migrate: ", err)
	}

	// Запучкаем мок-сервер music_info
	if config.MusicInfoAddr == "mock" {
		go mock.RunServer()
	}
	// Запускаем сервер приложения
	server.RunServer()
}
