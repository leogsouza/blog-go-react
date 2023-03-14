package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/leogsouza/blogbackend/database"
)

func main() {
	database.Connect()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	port := os.Getenv("PORT")
	app := fiber.New()
	app.Listen(":" + port)

}