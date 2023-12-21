package handler

import (
	"github.com/gofiber/fiber/v2"

	"invoicething/view/home"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h HomeHandler) HandleHomeShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("user") != nil
	return render(c, home.Show(isLoggedIn, c.Path()))
}
