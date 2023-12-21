package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setAuthCookies(c *fiber.Ctx, cookies AppCookies) {
	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: cookies.User,
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: cookies.Token,
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: cookies.RefreshToken,
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "expires_at",
		Value: strconv.FormatInt(cookies.ExpiresAt, 10),
		Path:  "/",
	})
}
