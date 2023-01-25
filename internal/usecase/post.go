package usecase

import (
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository"
	"github.com/v.kirpichov/db_tp/internal/utils/constants"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type PostUseCase struct {
	postRepository   repository.PostR
	forumRepository  ForumU
	userRepository   UserU
	threadRepository ThreadU
}

func NewPostUsecase(postRepository repository.PostR,
	forumRepository ForumU,
	userRepository UserU,
	threadRepository ThreadU) *PostUseCase {
	return &PostUseCase{
		postRepository:   postRepository,
		forumRepository:  forumRepository,
		userRepository:   userRepository,
		threadRepository: threadRepository,
	}
}

func (usecase *PostUseCase) Get(id int, details []string) (postDetailed *models.ParamsPost, err error) {
	postDetailed = &models.ParamsPost{}
	postDetailed.Post, err = usecase.postRepository.Get(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.PostNotFound
			return
		}
		err = errors.ServerInternal
		return
	}

	for _, detailType := range details {
		switch detailType {
		case constants.PostUser:
			postDetailed.Author, err = usecase.userRepository.Get(&postDetailed.Post.Author)
			if err != nil {
				postDetailed = nil
				return
			}
		case constants.PostThread:
			postDetailed.Thread, err = usecase.threadRepository.Get(strconv.Itoa(postDetailed.Post.Thread))
			if err != nil {
				postDetailed = nil
				return
			}
		case constants.PostForum:
			postDetailed.Forum, err = usecase.forumRepository.Get(postDetailed.Post.Forum)
			if err != nil {
				postDetailed = nil
				return
			}
		default:
			postDetailed = nil
			err = errors.BadRequest.SetTextDetails("неверные query параметры")
			return
		}
	}

	return
}

func (usecase *PostUseCase) Update(post *models.Post) (updatedPost *models.Post, err error) {
	updatedPost, err = usecase.postRepository.Update(post)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.PostNotFound
			return
		}
		err = errors.ServerInternal
		return
	}

	return
}
