package handler

import (
	"github.com/gofiber/fiber/v2"

	"invoicething/database"
	"invoicething/view/auth"

	authops "invoicething/ops/auth"
)

type AuthHandler struct {
	UserDB database.IDB
}

func NewAuthHandler(userDB database.IDB) *AuthHandler {
	return &AuthHandler{
		UserDB: userDB,
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

	err := h.UserDB.CreateUser(c.Context(), email, pwrd)
	if err != nil {
		return render(c, auth.SignUpForm(email, pwrd, "", "User already exists."))
	}

	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: email,
		Path:  "/",
	})

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

	e, err := authops.Login(c.Context(), h.UserDB, email, pwrd)
	if err != nil {
		var errMsg string
		switch err {
		case database.ErrUserNotFound:
			errMsg = "User not found."
		case database.ErrWrongPassword:
			errMsg = "Wrong password."
		default:
			errMsg = "Something went wrong."
		}

		return render(c, auth.LoginForm(email, pwrd, errMsg))
	}

	c.Cookie(&fiber.Cookie{
		Name:  "user",
		Value: e,
		Path:  "/",
	})
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
