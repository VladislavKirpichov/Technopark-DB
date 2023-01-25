package impl

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/queries"
)

type ForumRepository struct {
	db *pgxpool.Pool
}

func NewForumRepository(db *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (r *ForumRepository) Create(forum *models.Forum) (*models.Forum, error) {
	createdForum := &models.Forum{}

	row := r.db.QueryRow(context.Background(), queries.ForumQuery["Create"], forum.Slug, forum.Title, forum.User)

	err := row.Scan(
		&createdForum.Slug,
		&createdForum.Title,
		&createdForum.User,
		&createdForum.Posts,
		&createdForum.Threads,
	)

	if err != nil {
		return nil, err
	}

	return createdForum, nil
}

func (r *ForumRepository) Get(slug string) (*models.Forum, error) {
	row := r.db.QueryRow(context.Background(), queries.ForumQuery["Get"], slug)
	forum := &models.Forum{}

	err := row.Scan(
		&forum.ID,
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (repo *ForumRepository) CreateThread(thread *models.Thread) (*models.Thread, error) {
	row := repo.db.QueryRow(context.Background(), queries.ForumQuery["CreateThread"], thread.Slug, thread.Author, thread.Forum, thread.Title, thread.Msg, thread.Created)

	createdThread := &models.Thread{}
	err := row.Scan(
		&createdThread.ID,
		&createdThread.Slug,
		&createdThread.Author,
		&createdThread.Forum,
		&createdThread.Title,
		&createdThread.Msg,
		&createdThread.Created,
		&createdThread.Votes,
	)

	if err != nil {
		return nil, err
	}

	return createdThread, nil
}

func (r *ForumRepository) GetThreads(slug string, params *models.ForumQueryParams) ([]*models.Thread, error) {
	query := queries.ForumQuery["GetThreads"]
	var rows pgx.Rows
	var err error

	if !params.Since.Equal(time.Time{}) {
		if params.Desc {
			query = strings.Join([]string{query, queries.ForumQuery["GetThreadsDesc"]}, "")
		} else {
			query = strings.Join([]string{query, queries.ForumQuery["GetThreadsNoDesc"]}, "")
		}
		rows, err = r.db.Query(context.Background(), query, slug, params.Since, params.Limit)
	} else {
		if params.Desc {
			query = strings.Join([]string{query, queries.ForumQuery["GetThreadsSinceDesc"]}, "")
		} else {
			query = strings.Join([]string{query, queries.ForumQuery["GetThreadsSinceNoDesc"]}, "")
		}
		rows, err = r.db.Query(context.Background(), query, slug, params.Limit)
	}

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Slug,
			&thread.Author,
			&thread.Forum,
			&thread.Title,
			&thread.Msg,
			&thread.Created,
			&thread.Votes)

		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (r *ForumRepository) GetUsers(slug string, params *models.ForumUserQueryParams) ([]*models.User, error) {
	query := queries.ForumQuery["GetUsers"]
	var rows pgx.Rows
	var err error

	if params.Since != "" {
		if params.Desc {
			query = strings.Join([]string{query, queries.ForumQuery["GetUsersSinceDesc"]}, "")
		} else {
			query = strings.Join([]string{query, queries.ForumQuery["GetUsersSinceNoDesc"]}, "")
		}
		rows, err = r.db.Query(context.Background(), query, slug, params.Since, params.Limit)
	} else {
		if params.Desc {
			query = strings.Join([]string{query, queries.ForumQuery["GetUsersDesc"]}, "")
		} else {
			query = strings.Join([]string{query, queries.ForumQuery["GetUsersNoDesc"]}, "")
		}
		rows, err = r.db.Query(context.Background(), query, slug, params.Limit)
	}

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.Username,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
