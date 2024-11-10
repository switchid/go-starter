package handlers

import (
	appAuth "GoStarter/internal/handlers/app/auth"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Auth appAuth.AuthHandler
}

func New() *Handlers {
	return &Handlers{
		Auth: appAuth.AuthHandler{},
	}
}

func (h *Handlers) Index() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Apps",
		}, "layouts/main")
	}
}

func (h *Handlers) DashboardIndex() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("views/dashboard", fiber.Map{
			"Title": "Dashboard",
		}, "layouts/main")
	}
}

func (h *Handlers) NotFoundIndex() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("notfound", fiber.Map{
			"Title": "Not Found",
			"Data": fiber.Map{
				"Urls": c.Query("urls"),
			},
		}, "layouts/main")
	}
}
