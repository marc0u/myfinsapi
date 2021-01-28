package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/marc0u/myfinsapi/api/handlers"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {
	apiVersion := "2.2.0"
	fmt.Printf("---------- MyfinsAPI v%s ----------\n", apiVersion)
	var server = handlers.Server{}
	// Initialize DB
	server.InitializeDB(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	// Initialize Server
	server.RunServer(os.Getenv("API_PORT"), apiVersion)
}
