package handler

import (
	"github.com/gofiber/fiber/v2"

	"invoicething/database"
	"invoicething/external/supabase"
	"invoicething/view/auth"
)

type AuthHandler struct {
	userDB database.IDB
	sb     supabase.Client
}

func NewAuthHandler(userDB database.IDB, sb supabase.Client) *AuthHandler {
	return &AuthHandler{
		userDB: userDB,
		sb:     sb,
	}
}

func (h *AuthHandler) HandleSignupShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("logged_in").(bool)
	return render(c, auth.ShowSignup(isLoggedIn, c.Path()))
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

	setAuthCookies(c, email, res.AccessToken, res.RefreshToken)

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLoginShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("logged_in").(bool)
	return render(c, auth.ShowLogin(isLoggedIn, c.Path()))
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")

	res, err := h.sb.SigninUser(email, pwrd)
	if err != nil {
		return render(c, auth.LoginForm(email, pwrd, err.Error()))
	}

	setAuthCookies(c, res.User.Email, res.AccessToken, res.RefreshToken)

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: "",
		Path:  "/",
	})

	return c.Redirect("/")
}

func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	usr := c.Cookies("user")
	c.Locals("logged_in", usr != "")

	return c.Next()
}

func setAuthCookies(c *fiber.Ctx, email, token, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: email,
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/",
	})
}
