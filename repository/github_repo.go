package repository

import (
	"context"
	"github.com/nthduc/github-trending/models"
)

type GithubRepo interface {
	SaveRepo(context context.Context, user models.GithubRepo) (models.GithubRepo, error)
	SelectRepos(context context.Context, userId string, limit int) ([]models.GithubRepo, error)
	SelectRepoByName(context context.Context, name string) (models.GithubRepo, error)
	UpdateRepo(context context.Context, user models.GithubRepo) (models.GithubRepo, error)

	//Bookmark
	SelectAllBookmarks(context context.Context, userId string) ([]models.GithubRepo, error)
	Bookmark(context context.Context, bid, nameRepo, userId string) error
	DelBookmark(context context.Context, nameRepo, userId string) error
}
