package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

var sessionStore *session.Store

func init() {
	sessionConfig := session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:session_id",
		KeyGenerator: utils.UUIDv4,
	}
	sessionStore = session.New(sessionConfig)
}

func AppMiddleware(c *fiber.Ctx) error {
	storeSess, errSess := sessionStore.Get(c)
	if errSess != nil {
		fmt.Println(errSess)
	}

	//if storeSess.Get("login") == nil {
	//	storeSess.Set("login", false)
	//	err := storeSess.Save()
	//	if err != nil {
	//		return err
	//	}
	//}
	//fmt.Println("test")
	c.Locals("session", storeSess)

	err := c.Bind(fiber.Map{"AppName": "test"})
	if err != nil {
		return err
	}

	return c.Next()
}

func AuthMiddleware(c *fiber.Ctx) error {
	fmt.Println("auth middleware")
	sess := c.Locals("session").(*session.Session)

	if sess.Get("login") == nil {
		sess.Set("login", false)
	}

	logged := sess.Get("login").(bool)
	if logged == false {
		if c.Route().Name != "" {
			sess.Set("last_visit_route", c.Route().Name)
			err := sess.Save()
			if err != nil {
				return err
			}
		}
		err := c.RedirectToRoute("auth.login", fiber.Map{})
		if err != nil {
			return err
		}
	} else {
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}
	return c.Next()
}
