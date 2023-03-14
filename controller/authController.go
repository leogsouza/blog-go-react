package controller

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/leogsouza/blog-go-react/database"
	"github.com/leogsouza/blog-go-react/models"
)

func validateEmail(email string) bool {
	re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]`)

	return re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 character",
		})
	}

	email := strings.TrimSpace(data["email"].(string))

	if !validateEmail(email) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}

	database.DB.Where("email=?", email).First(&userData)
	if userData.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email already exist",
		})
	}
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     email,
	}

	user.SetPassword(data["password"].(string))

	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)

		return c.Status(500).JSON(fiber.Map{
			"message": "Cannot create the user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Account created successfully",
		"user":    user,
	})
}
