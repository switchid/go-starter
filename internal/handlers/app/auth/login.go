package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (a *AuthHandler) LoginIndex() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		return ctx.Render("views/auth/login", fiber.Map{
			"Title": "Login",
			"option": []map[string]string{
				{"value": "tes", "text": "test"},
				{"value": "tes2", "text": "test2"},
			},
		}, "layouts/auth")
	}
}
