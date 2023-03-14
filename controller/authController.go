package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/leogsouza/blog-go-react/models"
)

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
}
