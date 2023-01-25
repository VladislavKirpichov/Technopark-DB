package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/queries"
)

type ThreadRepository struct {
	db *pgxpool.Pool
}

func NewThreadRepository(db *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (repo *ThreadRepository) GetBySlug(slug string) (*models.Thread, error) {
	row := repo.db.QueryRow(context.Background(), queries.ThreadQuery["GetBySlug"], slug)

	thread := &models.Thread{}
	err := row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Msg, &thread.Created, &thread.Votes)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (repo *ThreadRepository) GetByID(id int) (*models.Thread, error) {
	row := repo.db.QueryRow(context.Background(), queries.ThreadQuery["GetByID"], id)

	thread := &models.Thread{}
	err := row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Msg, &thread.Created, &thread.Votes)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (repo *ThreadRepository) UpdateBySlug(thread *models.Thread) (*models.Thread, error) {
	row := repo.db.QueryRow(context.Background(), queries.ThreadQuery["UpdateBySlug"], thread.Title, thread.Msg, thread.Slug)
	updatedThread := &models.Thread{}

	err := row.Scan(&updatedThread.ID, &updatedThread.Slug, &updatedThread.Author, &updatedThread.Forum,
		&updatedThread.Title, &updatedThread.Msg, &updatedThread.Created, &updatedThread.Votes)

	if err != nil {
		return nil, err
	}

	return updatedThread, nil
}

func (repo *ThreadRepository) UpdateByID(thread *models.Thread) (*models.Thread, error) {
	row := repo.db.QueryRow(context.Background(), queries.ThreadQuery["UpdateByID"], thread.Title, thread.Msg, thread.ID)
	updatedThread := &models.Thread{}

	err := row.Scan(&updatedThread.ID, &updatedThread.Slug, &updatedThread.Author, &updatedThread.Forum,
		&updatedThread.Title, &updatedThread.Msg, &updatedThread.Created, &updatedThread.Votes)

	if err != nil {
		return nil, err
	}

	return updatedThread, nil
}

func (repo *ThreadRepository) VoteBySlug(slug string, vote *models.Vote) (err error) {
	_, err = repo.db.Exec(context.Background(), queries.ThreadQuery["VoteBySlug"], vote.Username, slug, vote.Voice)
	return
}

func (repo *ThreadRepository) VoteByID(id int, vote *models.Vote) (err error) {
	_, err = repo.db.Exec(context.Background(), queries.ThreadQuery["VoteByID"], vote.Username, id, vote.Voice)
	return
}

func (repo *ThreadRepository) CreatePostsBatch(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error) {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			trErr := tx.Rollback(ctx)
			if trErr != nil {
				err = trErr
			}
		} else {
			trErr := tx.Commit(ctx)
			if trErr != nil {
				err = trErr
			}
		}
	}()

	batch := new(pgx.Batch)
	createdTime := time.Now()

	for _, post := range posts {
		batch.Queue(queries.ThreadQuery["CreatePostsBatch"], post.Parent, post.Author, forumSlug, threadId, post.Message, createdTime)
	}

	batchRes := tx.SendBatch(ctx, batch)
	defer func() {
		batchErr := batchRes.Close()
		if batchErr != nil {
			err = batchErr
		}
	}()

	createdPosts := make([]*models.Post, 0)

	for i := 0; i < batch.Len(); i++ {
		createdPost := &models.Post{}

		row := batchRes.QueryRow()
		err := row.Scan(
			&createdPost.ID,
			&createdPost.Parent,
			&createdPost.Author,
			&createdPost.Forum,
			&createdPost.Thread,
			&createdPost.Created,
			&createdPost.IsEdited,
			&createdPost.Message)

		if err != nil {
			return nil, err
		}

		createdPosts = append(createdPosts, createdPost)
	}

	return createdPosts, nil
}

func (repo *ThreadRepository) CreatePosts(threadId int, forumSlug string, posts []*models.Post) ([]*models.Post, error) {
	query := queries.ThreadQuery["PostsCreate"]

	createdTime := time.Now()
	var values []interface{}
	for i, _ := range posts {
		indexShift := 6 * i
		query = strings.Join([]string{query, fmt.Sprintf("(NULLIF($%v, 0), $%v, $%v, $%v, $%v, $%v),",
			indexShift+1,
			indexShift+2,
			indexShift+3,
			indexShift+4,
			indexShift+5,
			indexShift+6)},
			"",
		)
		values = append(values, posts[i].Parent, posts[i].Author, forumSlug, threadId, posts[i].Message, createdTime)
	}
	query = strings.TrimSuffix(query, ",")
	query = strings.Join([]string{query, queries.ThreadQuery["CreatePostsTwo"]}, "")

	rows, err := repo.db.Query(context.Background(), query, values...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	createdPosts := make([]*models.Post, 0)
	for rows.Next() {
		createdPost := &models.Post{}
		err = rows.Scan(
			&createdPost.ID,
			&createdPost.Parent,
			&createdPost.Author,
			&createdPost.Forum,
			&createdPost.Thread,
			&createdPost.Created,
			&createdPost.IsEdited,
			&createdPost.Message)

		if err != nil {
			return nil, err
		}

		createdPosts = append(createdPosts, createdPost)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return createdPosts, nil
}

func (repo *ThreadRepository) GetPosts(threadId int, params *models.PostsQueryParams) ([]*models.Post, error) {
	var rows pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			rows, err = repo.db.Query(context.Background(), queries.DescNoSincePostQuery[params.SortType],
				threadId, params.Limit)
		} else {
			rows, err = repo.db.Query(context.Background(), queries.AscNoSincePostQuery[params.SortType],
				threadId, params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = repo.db.Query(context.Background(), queries.DescSincePostQuery[params.SortType],
				threadId, params.Since, params.Limit)
		} else {
			rows, err = repo.db.Query(context.Background(), queries.AscSincePostQuery[params.SortType],
				threadId, params.Since, params.Limit)
		}
	}

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.IsEdited,
			&post.Message)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
