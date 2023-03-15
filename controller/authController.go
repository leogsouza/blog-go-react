package controller

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leogsouza/blog-go-react/database"
	"github.com/leogsouza/blog-go-react/models"
	"github.com/leogsouza/blog-go-react/util"
	"gorm.io/gorm"
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

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User

	if err := db.Where(&models.User{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return &user, nil
}

func Login(c *fiber.Ctx) error {

	type loginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(loginInput)
	var userData *models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error on login request", "data": err,
		})
	}

	userData, err := getUserByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Error on login request", "data": err,
		})
	}

	if err := userData.ComparePassword(input.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials", "data": err,
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(userData.ID)))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials", "data": err,
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successfully login",
		"user":    userData,
	})

}
