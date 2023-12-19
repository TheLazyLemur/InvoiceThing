package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"invoicething/database"
	"invoicething/handler"
)

func main() {
	app := fiber.New()

	udb := database.NewLiteDB()

	homeHandler := handler.HomeHandler{}
	authHandler := handler.NewAuthHandler(udb)
	dashboardHanlder := handler.NewDashboardHandler(udb)

	app.Use(authHandler.AuthMiddleware)

	app.Get("/", homeHandler.HandleHomeShow)

	app.Get("/auth/signup", authHandler.HandleSignupShow)
	app.Post("/auth/signup", authHandler.HandleSignup)
	app.Get("/auth/login", authHandler.HandleLoginShow)
	app.Post("/auth/login", authHandler.HandleLogin)
	app.Get("/auth/logout", authHandler.HandleLogout)
	app.Get("/dashboard", dashboardHanlder.HandleShowDashboard)

	fmt.Println("Listening on port 3000")
	app.Listen(":3000")
}
