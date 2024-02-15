package repository

import (
	"context"

	"github.com/nthduc/github-trending/models"
	"github.com/nthduc/github-trending/models/req"
)

type UserRepo interface {
	CheckLogin(context context.Context, loginReq req.ReqSignIn) (models.User, error)
	SaveUser(context context.Context, user models.User) (models.User, error)
	SelectUserById(context context.Context, userId string) (models.User, error)
	UpdateUser(context context.Context, user models.User) (models.User, error)
}
