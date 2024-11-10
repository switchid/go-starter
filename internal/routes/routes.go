package routes

import (
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	App *fiber.App
}

func Load(app *fiber.App) *Routes {
	return &Routes{app}
}
