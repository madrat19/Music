package tools

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	DBName        string
	Username      string
	Password      string
	ServerAddr    string
	MusicInfoAddr string
}

var config *Config

// Получает настройки
func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		err := godotenv.Load("../.env")
		if err != nil {
			panic("no .env file provided!")
		}
		config.Host = os.Getenv("PGHOST")
		config.Port = os.Getenv("PGPORT")
		config.DBName = os.Getenv("DBNAME")
		config.Username = os.Getenv("PGUSER")
		config.Password = os.Getenv("PGPASSWORD")
		config.ServerAddr = os.Getenv("SERVER")
		config.MusicInfoAddr = os.Getenv("MUSICINFO")
	}

	return config
}
