package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Rosya-edwica/go-trudvsem/api"
	"github.com/Rosya-edwica/go-trudvsem/db"
	"github.com/joho/godotenv"
)


func initDatabase() (database *db.Database) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Создайте файл с переменными окружения .env")
	}

	database = &db.Database{
		Host: os.Getenv("MYSQL_HOST"),
		Port: os.Getenv("MYSQL_PORT"),
		User: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Name: os.Getenv("MYSQL_DATABASE"),
	}
	return database
}

func main() {
	startTime := time.Now().Unix()
	database := initDatabase()
	database.Connect()
	defer database.Close()

	cities := database.GetCities()
	positions := database.GetPositions() 
	for _, item := range positions {
		fmt.Println(item.Name)
		vacancies := api.FindVacanciesByPosition(item, cities)
		database.SaveManyVacancies(vacancies)
	}

	fmt.Println("Время выполнения: ", time.Now().Unix() - startTime)
}