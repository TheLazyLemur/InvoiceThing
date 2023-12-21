package main

import (
	"log/slog"
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

	app.Get("/auth/signup", authHandler.OptionallyProtectedRoute, authHandler.HandleSignupShow)
	app.Post("/auth/signup", authHandler.OptionallyProtectedRoute, authHandler.HandleSignup)
	app.Get("/auth/login", authHandler.OptionallyProtectedRoute, authHandler.HandleLoginShow)
	app.Post("/auth/login", authHandler.OptionallyProtectedRoute, authHandler.HandleLogin)
	app.Get("/auth/logout", authHandler.OptionallyProtectedRoute, authHandler.HandleLogout)

	app.Get("/", authHandler.OptionallyProtectedRoute, homeHandler.HandleHomeShow)

	app.Get("/dashboard", authHandler.ProtectedRoute, dashboardHanlder.HandleShowDashboard)

	slog.Info("Listening on port 3000")
	app.Listen(":3000")
}
