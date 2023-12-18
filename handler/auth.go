package handler

import (
	"fmt"

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
	return render(c, auth.ShowSignup(isLoggedIn))
}

func (h *AuthHandler) HandleSignup(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")
	cpwrd := c.FormValue("confirm_password")

	if pwrd != cpwrd {
		return render(c, auth.SignUpForm([]string{"Passwords do not match."}, email, pwrd, ""))
	}

	err := h.UserDB.CreateUser(c.Context(), email, pwrd)
	if err != nil {
		return render(c, auth.SignUpForm([]string{"User already exists."}, email, pwrd, ""))
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
	return render(c, auth.ShowLogin(isLoggedIn))
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pwrd := c.FormValue("password")

	e, err := authops.Login(c.Context(), h.UserDB, email, pwrd)
	if err != nil {
		errs := []string{}

		switch err {
		case database.ErrUserNotFound:
			errs = append(errs, "User not found.")
		case database.ErrWrongPassword:
			errs = append(errs, "Wrong password.")
		default:
			errs = append(errs, "Something went wrong.")
		}

		return render(c, auth.LoginForm(errs, email, pwrd))
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
	if usr == "" {
		fmt.Println("Not Logged in")
		c.Locals("logged_in", false)
	} else {
		fmt.Println("Logged in")
		c.Locals("logged_in", true)
	}

	return c.Next()
}
