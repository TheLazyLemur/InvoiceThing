package handler

import (
	"invoicething/view/dashboard"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) HandleShowDashboard(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("user") != nil
	return render(c, dashboard.Show(isLoggedIn, c.Path()))
}
