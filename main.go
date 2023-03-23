package main

import (
	"assignment2/api/controllers"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env file")
	}
}

func main() {
	fmt.Println("Starting assignment2")
	fmt.Println("running server")

	loadEnvFile()
	app := controllers.App{}
	app.Initialize()
	app.RunServer()

}
