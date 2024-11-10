package routes

import (
	"GoStarter/internal/handlers"
	"GoStarter/internal/middleware"
)

func (r *Routes) AppRoutes() {
	handler := handlers.New()
	r.App.Get("/", handler.Index()).Name("index")
	r.App.Get("/dashboard", middleware.AuthMiddleware, handler.DashboardIndex()).Name("dashboard")

	//authGroup := r.App.Group("auth")
	//authGroup.Get("login", handler.Auth.Login())
}

func (r *Routes) AuthRoutes() {
	handler := handlers.New()
	g := r.App.Group("/auth")
	g.Get("/login", handler.Auth.LoginIndex()).Name("auth.login.index")
	g.Post("/login", handler.Auth.Login()).Name("auth.login")
	g.Get("/register", handler.Auth.RegisterIndex()).Name("auth.register.index")
	g.Post("/register", handler.Auth.Register()).Name("auth.register")
	g.Get("/logout", handler.Auth.Logout()).Name("auth.logout")
}

func (r *Routes) NotFoundRoutes() {
	handler := handlers.New()
	r.App.Get("/not-found", handler.NotFoundIndex()).Name("not-found")
}
