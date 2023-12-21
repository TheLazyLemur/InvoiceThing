package handler

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"invoicething/external/supabase"
	"invoicething/view/auth"
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

	setAuthCookies(c, email, res.AccessToken, res.RefreshToken, res.ExpiresAt)

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLoginShow(c *fiber.Ctx) error {
	isLoggedIn := c.Locals("user") != nil
	return render(c, auth.ShowLogin(isLoggedIn, c.Path()))
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")

	res, err := h.sb.SigninUser(email, pwrd)
	if err != nil {
		return render(c, auth.LoginForm(email, pwrd, err.Error()))
	}

	setAuthCookies(c, res.User.Email, res.AccessToken, res.RefreshToken, res.ExpiresAt)

	c.Response().Header.Set("HX-Redirect", "/")
	return nil
}

func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: "",
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: "",
		Path:  "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "expires_at",
		Path:  "/",
		Value: "0",
	})

	return c.Redirect("/")
}

func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("token")
	expiresAt, parseError := strconv.ParseInt(c.Cookies("expires_at"), 10, 64)
	if parseError != nil {
		expiresAt = 0
		return c.Next()
	}
	refreshToken := c.Cookies("refresh_token")

	var user supabase.User
	var err error

	currentTime := time.Now().Unix()
	if expiresAt < currentTime {
		slog.Info("Attempting to refresh token")

		resp, err := h.sb.RefreshToken(refreshToken)
		if err != nil {
			return c.Next()
		}

		user = resp.User
		setAuthCookies(c, resp.User.Email, resp.AccessToken, resp.RefreshToken, resp.ExpiresAt)
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

func setAuthCookies(c *fiber.Ctx, email string, token string, refreshToken string, expiresAt int64) {
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

	c.Cookie(&fiber.Cookie{
		Name:  "expires_at",
		Path:  "/",
		Value: strconv.FormatInt(expiresAt, 10),
	})
}
