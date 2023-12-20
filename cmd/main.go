package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"invoicething/database"
	"invoicething/external/supabase"
	"invoicething/handler"
)

func main() {
	godotenv.Load()
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	app := fiber.New()

	udb := database.NewLiteDB()
	sb := supabase.NewClient(supabaseURL, supabaseKey)

	homeHandler := handler.HomeHandler{}
	authHandler := handler.NewAuthHandler(udb, sb)
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
