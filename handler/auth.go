package handler

import (
	"invoicething/external/supabase"
	"invoicething/view/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	sb supabase.Client
}

func NewAuthHandler(sb supabase.Client) *AuthHandler {
	return &AuthHandler{
		sb: sb,
	}
}

func (h *AuthHandler) HandleSignupShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("user") != nil
	return render(c, auth.ShowSignup(isLoggedIn, c.Path(), "https://avatars.githubusercontent.com/u/36530232?v=4"))
}

func (h *AuthHandler) HandleSignup(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")
	cpwrd := c.FormValue("confirm_password")

	if pwrd != cpwrd {
		return render(c, auth.SignUpForm(email, pwrd, "", "Passwords do not match."))
	}

	res, err := h.sb.SignupUser(email, pwrd)
	if err != nil {
		return render(c, auth.SignUpForm(email, pwrd, "", err.Error()))
	}

	setAuthCookies(c, AppCookies{
		User:         email,
		Token:        res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    res.ExpiresAt,
		Path:         "/",
	})

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLoginShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("user") != nil
	return render(c, auth.ShowLogin(isLoggedIn, c.Path(), "https://avatars.githubusercontent.com/u/36530232?v=4"))
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")

	res, err := h.sb.SigninUser(email, pwrd)
	if err != nil {
		return render(c, auth.LoginForm(email, pwrd, err.Error()))
	}

	setAuthCookies(c, AppCookies{
		User:         email,
		Token:        res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    res.ExpiresAt,
		Path:         "/",
	})

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	setAuthCookies(c, AppCookies{})

	return c.Redirect("/")
}
