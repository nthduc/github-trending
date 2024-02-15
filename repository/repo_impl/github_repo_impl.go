package repo_impl

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/nthduc/github-trending/db"
	"github.com/nthduc/github-trending/exceptions"
	"github.com/nthduc/github-trending/log"
	"github.com/nthduc/github-trending/models"
	"github.com/nthduc/github-trending/repository"
)

type GithubRepoImpl struct {
	sql *db.Sql
}

func NewGithubRepo(sql *db.Sql) repository.GithubRepo {
	return &GithubRepoImpl{
		sql: sql,
	}
}

func (g GithubRepoImpl) SelectRepoByName(context context.Context, name string) (models.GithubRepo, error) {
	var repo = models.GithubRepo{}
	err := g.sql.Db.GetContext(context, &repo,
		`SELECT * FROM repos WHERE name=$1`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return repo, exceptions.RepoNotFound
		}
		log.Error(err.Error())
		return repo, err
	}
	return repo, nil
}

func (g GithubRepoImpl) SaveRepo(context context.Context, repo models.GithubRepo) (models.GithubRepo, error) {
	// name, description, url, color, lang, fork, stars, stars_today, build_by, created_at, updated_at
	statement := `INSERT INTO repos(
					name, description, url, color, lang, fork, stars, 
 			        stars_today, build_by, created_at, updated_at) 
          		  VALUES(
					:name,:description, :url, :color, :lang, :fork, :stars, 
					:stars_today, :build_by, :created_at, :updated_at
				  )`

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err := g.sql.Db.NamedExecContext(context, statement, repo)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return repo, exceptions.RepoConflict
			}
		}
		log.Error(err.Error())
		return repo, exceptions.RepoInsertFail
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectRepos(context context.Context, userId string, limit int) ([]models.GithubRepo, error) {
	var repos []models.GithubRepo
	err := g.sql.Db.SelectContext(context, &repos,
		`
			SELECT 
				repos.name, repos.description, repos.url, repos.color, repos.lang, 
				repos.fork, repos.stars, repos.stars_today, repos.build_by, repos.updated_at, 
				COALESCE(repos.name = bookmarks.repo_name, FALSE) as bookmarked
			FROM repos
			FULL OUTER JOIN bookmarks 
			ON repos.name = bookmarks.repo_name AND 
			   bookmarks.user_id=$1  
			WHERE repos.name IS NOT NULL 
			ORDER BY updated_at ASC LIMIT $2
		`, userId, limit)

	if err != nil {
		log.Error(err.Error())
		return repos, err
	}

	return repos, nil
}

func (g GithubRepoImpl) UpdateRepo(context context.Context, repo models.GithubRepo) (models.GithubRepo, error) {
	// name, description, url, color, lang, fork, stars, stars_today, build_by, created_at, updated_at
	sqlStatement := `
		UPDATE repos
		SET 
			stars  = :stars,
			fork = :fork,
			stars_today = :stars_today,
			build_by = :build_by,
			updated_at = :updated_at
		WHERE name = :name
	`

	result, err := g.sql.Db.NamedExecContext(context, sqlStatement, repo)
	if err != nil {
		log.Error(err.Error())
		return repo, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return repo, exceptions.RepoNotUpdated
	}
	if count == 0 {
		return repo, exceptions.RepoNotUpdated
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectAllBookmarks(context context.Context, userId string) ([]models.GithubRepo, error) {
	repos := []models.GithubRepo{}
	err := g.sql.Db.SelectContext(context, &repos,
		`SELECT 
					repos.name, repos.description, repos.url, 
					repos.color, repos.lang, repos.fork, repos.stars, 
					repos.stars_today, repos.build_by, true as bookmarked
				FROM bookmarks 
				INNER JOIN repos
				ON bookmarks.user_id=$1 AND repos.name = bookmarks.repo_name`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return repos, exceptions.BookmarkNotFound
		}
		log.Error(err.Error())
		return repos, err
	}
	return repos, nil
}

func (g GithubRepoImpl) Bookmark(context context.Context, bid, nameRepo, userId string) error {
	statement := `INSERT INTO bookmarks(
					bid, user_id, repo_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`

	now := time.Now()
	_, err := g.sql.Db.ExecContext(
		context, statement, bid, userId,
		nameRepo, now, now)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return exceptions.BookmarkConflic
			}
		}
		log.Error(err.Error())
		return exceptions.BookmarkFail
	}

	return nil
}

func (g GithubRepoImpl) DelBookmark(context context.Context, nameRepo, userId string) error {
	result := g.sql.Db.MustExecContext(
		context,
		"DELETE FROM bookmarks WHERE repo_name = $1 AND user_id = $2",
		nameRepo, userId)

	_, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return exceptions.DelBookmarkFail
	}

	return nil
}
