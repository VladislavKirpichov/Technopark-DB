package impl

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/queries"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Get(id int) (*models.Post, error) {
	row := r.db.QueryRow(context.Background(), queries.PostQuery["Get"], id)
	post := &models.Post{}
	err := row.Scan(
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

	return post, nil
}

func (r *PostRepository) Update(post *models.Post) (*models.Post, error) {
	row := r.db.QueryRow(context.Background(), queries.PostQuery["Update"], post.Message, post.ID)

	updatedPost := &models.Post{}
	err := row.Scan(
		&updatedPost.ID,
		&updatedPost.Parent,
		&updatedPost.Author,
		&updatedPost.Forum,
		&updatedPost.Thread,
		&updatedPost.Created,
		&updatedPost.IsEdited,
		&updatedPost.Message)

	if err != nil {
		return nil, err
	}
	
	return updatedPost, nil
}

