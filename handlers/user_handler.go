package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleSignIn(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"user":  "nthduc",
		"email": "nthduc@gmail.com",
	})
}

func HandleSignUp(c echo.Context) error {
	type User struct {
		Email    string
		FullName string
		Age      int
	}

	user := User{
		Email:    "nthduc@gmail.com",
		FullName: "nthduc",
		Age:      18,
	}
	return c.JSON(http.StatusOK, user)
}
