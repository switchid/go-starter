package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (a *AuthHandler) RegisterIndex() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("views/auth/register", fiber.Map{
			"Title": "Test Hello World",
		}, "layouts/main")
	}
}
