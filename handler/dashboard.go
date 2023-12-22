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
	return render(c, dashboard.Show(true, c.Path(), "https://avatars.githubusercontent.com/u/36530232?v=4"))
}
