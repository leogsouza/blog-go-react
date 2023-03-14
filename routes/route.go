package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leogsouza/blog-go-react/controller"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
}
