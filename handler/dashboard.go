package handler

import (
	"invoicething/database"
	"invoicething/view/dashboard"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
}

func NewDashboardHandler(userDB database.IDB) *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) HandleShowDashboard(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("logged_in").(bool)
	return render(c, dashboard.Show(isLoggedIn, c.Path()))
}
