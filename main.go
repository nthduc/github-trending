package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nthduc/github-trending/handlers"
)

func main() {
	e := echo.New()
	e.GET("/", welcome)

	e.GET("/user/sign-in", handlers.HandleSignIn)
	e.GET("/user/sign-up", handlers.HandleSignUp)
	e.Logger.Fatal(e.Start(":3000"))
}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Go")
}
