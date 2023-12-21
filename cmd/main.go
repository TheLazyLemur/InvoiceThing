package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"invoicething/external/supabase"
	"invoicething/handler"
)

func main() {
	godotenv.Load()
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	app := fiber.New()

	sb := supabase.NewClient(supabaseURL, supabaseKey)

	homeHandler := handler.NewHomeHandler()
	authHandler := handler.NewAuthHandler(sb)
	dashboardHanlder := handler.NewDashboardHandler()

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
