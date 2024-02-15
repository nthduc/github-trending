package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nthduc/github-trending/handlers"
	"github.com/nthduc/github-trending/middleware"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handlers.UserHandler
	RepoHandler handlers.RepoHandler
}

func (api *API) SetupRouter() {
	// user
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignIn)
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignUp)

	// profile
	user := api.Echo.Group("/user", middleware.JWTMiddleware())
	user.GET("/profile", api.UserHandler.Profile)
	user.PUT("/profile/update", api.UserHandler.UpdateProfile)

	// github repo
	github := api.Echo.Group("/github", middleware.JWTMiddleware())
	github.GET("/trending", api.RepoHandler.RepoTrending)

	// bookmark
	bookmark := api.Echo.Group("/bookmark", middleware.JWTMiddleware())
	bookmark.GET("/list", api.RepoHandler.SelectBookmarks)
	bookmark.POST("/add", api.RepoHandler.Bookmark)
	bookmark.DELETE("/delete", api.RepoHandler.DelBookmark)
}
