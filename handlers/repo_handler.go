package handlers

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nthduc/github-trending/log"
	"github.com/nthduc/github-trending/models"
	"github.com/nthduc/github-trending/models/req"
	"github.com/nthduc/github-trending/repository"
)

type RepoHandler struct {
	GithubRepo repository.GithubRepo
}

func (r RepoHandler) RepoTrending(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JwtCustomClaims)

	repos, _ := r.GithubRepo.SelectRepos(c.Request().Context(), claims.UserId, 25)
	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

func (r RepoHandler) SelectBookmarks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JwtCustomClaims)

	repos, _ := r.GithubRepo.SelectAllBookmarks(
		c.Request().Context(),
		claims.UserId)

	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

func (r RepoHandler) Bookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JwtCustomClaims)

	bId, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, models.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	err = r.GithubRepo.Bookmark(
		c.Request().Context(),
		bId.String(),
		req.RepoName,
		claims.UserId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark thành công",
		Data:       nil,
	})
}

func (r RepoHandler) DelBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*models.JwtCustomClaims)

	err = r.GithubRepo.DelBookmark(
		c.Request().Context(),
		req.RepoName, claims.UserId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Message:    "Xoá bookmark thành công",
		Data:       nil,
	})
}
