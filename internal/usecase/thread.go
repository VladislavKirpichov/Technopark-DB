package usecase

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository"
	"github.com/v.kirpichov/db_tp/internal/utils/check"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type ThreadUseCase struct {
	threadRepository repository.ThreadR
}

func NewThreadUseCase(threadRepository repository.ThreadR) *ThreadUseCase {
	return &ThreadUseCase{threadRepository: threadRepository}
}

func (usecase *ThreadUseCase) Get(slugOrId string) (thread *models.Thread, err error) {
	v, _ := check.GetInstance()
	slug, id, err := v.GetSlugOrIdOrErr(slugOrId)
	if err != nil {
		err = errors.BadRequest.SetTextDetails(err.Error())
		return
	}

	if slug == "" {
		thread, err = usecase.threadRepository.GetByID(id)
	} else {
		thread, err = usecase.threadRepository.GetBySlug(slug)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ThreadNotFound
		} else {
			err = errors.ServerInternal
		}
	}

	return
}

func (usecase *ThreadUseCase) Update(slugOrId string, thread *models.Thread) (updatedThread *models.Thread, err error) {
	v, _ := check.GetInstance()
	slug, id, err := v.GetSlugOrIdOrErr(slugOrId)
	if err != nil {
		err = errors.BadRequest.SetTextDetails(err.Error())
		return
	}

	if slug == "" {
		thread.ID = id
		updatedThread, err = usecase.threadRepository.UpdateByID(thread)
	} else {
		thread.Slug = slug
		updatedThread, err = usecase.threadRepository.UpdateBySlug(thread)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ThreadUpdateNotFound
			return
		}
		err = errors.ServerInternal
		return
	}

	return
}

func (usecase *ThreadUseCase) Vote(slugOrId string, vote *models.Vote) (thread *models.Thread, err error) {
	v, _ := check.GetInstance()
	if !v.CheckVote(vote) {
		err = errors.BadRequest.SetTextDetails("не верное значение голоса")
		return
	}

	slug, id, err := v.GetSlugOrIdOrErr(slugOrId)
	if err != nil {
		err = errors.BadRequest.SetTextDetails(err.Error())
		return
	}

	if slug == "" {
		err = usecase.threadRepository.VoteByID(id, vote)
	} else {
		err = usecase.threadRepository.VoteBySlug(slug, vote)
	}

	if err != nil {
		pgconErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgconErr.SQLState() == errors.Err23503 || pgconErr.SQLState() == errors.Err23502 {
				err = errors.ThreadUserOrThreadNotFound
				return
			} else {
				err = errors.ServerInternal.SetTextDetails(err.Error())
				return
			}
		}
	}

	if slug == "" {
		thread, err = usecase.threadRepository.GetByID(id)
	} else {
		thread, err = usecase.threadRepository.GetBySlug(slug)
	}

	if err != nil {
		err = errors.ServerInternal
		return
	}

	return
}

func (usecase *ThreadUseCase) CreatePosts(slugOrId string, posts []*models.Post) (createdPosts []*models.Post, err error) {
	thread, err := usecase.Get(slugOrId)
	if err != nil {
		return
	}

	if len(posts) == 0 {
		createdPosts = make([]*models.Post, 0)
		return
	}

	createdPosts, err = usecase.threadRepository.CreatePosts(thread.ID, thread.Forum, posts)
	if err != nil {
		pgconErr, ok := err.(*pgconn.PgError)
		if ok {
			switch pgconErr.SQLState() {
			case errors.Err23503:
				err = errors.PostUserNotFound
				return

			case errors.P0001:
				err = errors.PostWrongParent
				return

			default:
				err = errors.ServerInternal
			}
		}
	}

	return
}

func (usecase *ThreadUseCase) GetPosts(slugOrId string, params *models.PostsQueryParams) (posts []*models.Post, err error) {
	thread, err := usecase.Get(slugOrId)
	if err != nil {
		return
	}

	posts, err = usecase.threadRepository.GetPosts(thread.ID, params)
	if err != nil {
		err = errors.ServerInternal
	}

	return
}
