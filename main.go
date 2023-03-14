package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/leogsouza/blog-go-react/database"
	"github.com/leogsouza/blog-go-react/routes"
)

func main() {
	database.Connect()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	port := os.Getenv("PORT")
	log.Println(port)
	app := fiber.New()
	routes.Setup(app)
	log.Fatal(app.Listen(":" + port))

}
