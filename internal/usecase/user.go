package usecase

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type UserUsecase struct {
	userRepository repository.UserR
}

func NewUserUsecase(userRepository repository.UserR) *UserUsecase {
	return &UserUsecase{userRepository: userRepository}
}

func (usecase *UserUsecase) Get(nickname *string) (user *models.User, err error) {
	user, err = usecase.userRepository.Get(nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.NotFoundUser
		} else {
			err = errors.ServerInternal
		}
	}
	return
}

func (usecase *UserUsecase) All() (users *[]models.User, err error) {
	return
}

func (usecase *UserUsecase) Create(user *models.User) (users []*models.User, err error) {
	err = usecase.userRepository.Create(user)

	if err != nil {
		pgconErr, ok := err.(*pgconn.PgError)
		if ok && pgconErr.SQLState() == errors.Err23505 {
			users, err = usecase.userRepository.GetUsersByUserNicknameOrEmail(user)
			if err != nil {
				err = errors.ServerInternal
				return
			}
			err = errors.ConflictUserCreate
			return
		}
		err = errors.ServerInternal
		return
	}

	return
}

func (usecase *UserUsecase) Update(user *models.User) (updatedUser *models.User, err error) {
	updatedUser, err = usecase.userRepository.Update(user)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.NotFoundUserUpdate
			return
		}
		pgconErr, ok := err.(*pgconn.PgError)
		if ok && pgconErr.SQLState() == errors.Err23505 {
			err = errors.ConflictUserUpdate
			return
		}
		err = errors.ServerInternal
		return
	}

	return
}
