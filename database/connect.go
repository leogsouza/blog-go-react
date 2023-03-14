package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/leogsouza/blog-go-react/models"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	dsn := os.Getenv("DB_DSN")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}
	log.Println("connect successfully")

	DB = database

	database.AutoMigrate(
		&models.User{},
	)

}
