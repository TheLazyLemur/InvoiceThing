package handler

import (
	"invoicething/external/supabase"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) OptionallyProtectedRoute(c *fiber.Ctx) error {
	token := c.Cookies("token")
	refreshToken := c.Cookies("refresh_token")
	expiresAt, parseError := strconv.ParseInt(c.Cookies("expires_at"), 10, 64)
	if parseError != nil {
		expiresAt = 0
	}

	var user supabase.User
	var err error

	slog.Info("Checking auth optionally")
	currentTime := time.Now().Unix()

	if expiresAt < currentTime {
		if refreshToken == "" {
			return c.Next()
		}

		slog.Info("Attempting to refresh token")

		resp, err := h.sb.RefreshToken(refreshToken)
		if err != nil {
			slog.Error("Failed to refresh token", err)
			setAuthCookies(c, AppCookies{})
			return c.Next()
		}

		user = resp.User
		setAuthCookies(c, AppCookies{
			User:         resp.User.Email,
			Token:        resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
			Path:         "/",
		})
		slog.Info("Successfully refreshed token")
	} else {
		user, err = h.sb.GetUser(token)
		if err != nil {
			return c.Next()
		}
	}

	c.Locals("user", user)
	c.Locals("email", user.Email)
	c.Locals("token", token)
	c.Locals("refresh_token", refreshToken)

	return c.Next()
}

func (h *AuthHandler) ProtectedRoute(c *fiber.Ctx) error {
	token := c.Cookies("token")
	refreshToken := c.Cookies("refresh_token")
	expiresAt, parseError := strconv.ParseInt(c.Cookies("expires_at"), 10, 64)
	if parseError != nil {
		expiresAt = 0
	}

	var user supabase.User
	var err error

	slog.Info("Checking auth")
	currentTime := time.Now().Unix()

	if currentTime > expiresAt && refreshToken == "" {
		c.Redirect("/auth/login")
		return c.Next()
	}

	if expiresAt < currentTime {
		slog.Info("Attempting to refresh token")

		resp, err := h.sb.RefreshToken(refreshToken)
		if err != nil {
			slog.Error("Failed to refresh token", err)
			setAuthCookies(c, AppCookies{})
			return c.Next()
		}

		user = resp.User
		setAuthCookies(c, AppCookies{
			User:         resp.User.Email,
			Token:        resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
			Path:         "/",
		})
		slog.Info("Successfully refreshed token")
	} else {
		user, err = h.sb.GetUser(token)
		if err != nil {
			return c.Next()
		}
	}

	c.Locals("user", user)
	c.Locals("email", user.Email)
	c.Locals("token", token)
	c.Locals("refresh_token", refreshToken)

	return c.Next()
}
