package handler

import (
	"github.com/gofiber/fiber/v2"

	"invoicething/view/home"
)

type HomeHandler struct{}

func (h HomeHandler) HandleHomeShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("logged_in").(bool)
	return render(c, home.Show(isLoggedIn, c.Path()))
}
