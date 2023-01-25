package impl

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/utils/queries"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(nickname *string) (*models.User, error) {
	user := &models.User{}
	row := r.db.QueryRow(context.Background(), queries.UserQuery["Get"], *nickname)
	err := row.Scan(
		&user.Username,
		&user.FullName,
		&user.About,
		&user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) All() (*[]models.User, error) {
	return nil, nil
}

func (r *UserRepository) Create(user *models.User) (err error) {
	_, err = r.db.Exec(context.Background(), queries.UserQuery["Create"], user.Username, user.FullName, user.About, user.Email)
	return err
}

func (r *UserRepository) Update(user *models.User) (*models.User, error) {
	row := r.db.QueryRow(context.Background(), queries.UserQuery["Update"], user.FullName, user.About, user.Email, user.Username)
	updatedUser := &models.User{}
	err := row.Scan(
		&updatedUser.Username,
		&updatedUser.FullName,
		&updatedUser.About,
		&updatedUser.Email)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) GetUsersByUserNicknameOrEmail(user *models.User) ([]*models.User, error) {
	rows, err := r.db.Query(context.Background(), queries.UserQuery["GetUsersByUserNOE"], user.Username, user.Email)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)

	for rows.Next() {
		conflictUser := &models.User{}
		err = rows.Scan(
			&conflictUser.Username,
			&conflictUser.FullName,
			&conflictUser.About,
			&conflictUser.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, conflictUser)
	}

	return users, nil
}
