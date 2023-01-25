package impl

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/queries"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (repo *ServiceRepository) Clear() (err error) {
	_, err = repo.db.Exec(context.Background(), queries.ServiceQuery["Clear"])
	return err
}

func (repo *ServiceRepository) Status() (*models.ForumStatus, error) {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			trErr := tx.Commit(ctx)
			if trErr != nil {
				err = trErr
			}
		} else {
			trErr := tx.Rollback(ctx)
			if trErr != nil {
				err = trErr
			}
		}
	}()

	status := &models.ForumStatus{}
	if err = tx.QueryRow(ctx, queries.ServiceQuery["queryUsers"]).Scan(&status.User); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(ctx, queries.ServiceQuery["queryForums"]).Scan(&status.Forum); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(ctx, queries.ServiceQuery["queryThreads"]).Scan(&status.Thread); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(ctx, queries.ServiceQuery["queryPosts"]).Scan(&status.Post); err != nil {
		return nil, err
	}

	return status, nil
}
