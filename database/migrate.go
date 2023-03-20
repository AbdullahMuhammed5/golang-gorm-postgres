package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/initializers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/utils"
)

func init() {
	initializers.InitializeEnv(".", "app")
}

func runManualMigration() {
	// get files
	files, err := utils.IOReadDir("./database/migrations")
	if err != nil {
		panic("Error while reading migrations directory.")
	}
	// loop and excute
	for _, file := range files {
		query, err := ioutil.ReadFile(path.Join("./database/migrations", file))
		if err != nil {
			panic(err)
		}
		if err := initializers.AppInstance.DB.Exec(string(query[:])).Error; err != nil {
			panic(err)
		}
	}
}

func main() {
	runManualMigration()
	initializers.AppInstance.DB.AutoMigrate(&models.User{})
	initializers.AppInstance.DB.AutoMigrate(&models.Ticket{})
	fmt.Println("Migration complete")
}
