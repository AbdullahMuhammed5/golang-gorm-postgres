package initializers

import (
	"log"

	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Config *Config
}

var AppInstance App

func InitializeEnv(path string, filename string) App {
	config, err := LoadConfig(path, filename)
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	db := ConnectDB(&config)
	AppInstance = App{DB: db, Config: &config}
	return App{DB: db, Config: &config}
}
