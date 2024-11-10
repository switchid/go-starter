package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthHandler struct{}

func (a *AuthHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess := c.Locals("session").(*session.Session)
		lastRoute := ""
		if sess.Get("last_visit_route") != nil {
			lastRoute = sess.Get("last_visit_route").(string)
		}

		sess.Set("login", true)
		if err := sess.Save(); err != nil {
			panic(err)
		}

		return c.RedirectToRoute(lastRoute, fiber.Map{})
	}
}

func (a *AuthHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.RedirectToRoute("index", fiber.Map{})
	}
}

func (a *AuthHandler) Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess := c.Locals("session").(*session.Session)

		if sess.Get("login").(bool) {
			sess.Set("login", false)
			err := sess.Save()
			if err != nil {
				return err
			}
		}

		return c.RedirectToRoute("auth.login", fiber.Map{})
	}
}
