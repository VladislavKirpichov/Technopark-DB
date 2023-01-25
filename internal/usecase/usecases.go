package usecase

import (
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository"
)

type Usecases struct {
	ForumUsecase   ForumU
	PostUsecase    PostU
	UserUsecase    UserU
	ServiceUsecase ServiceU
	ThreadUsecase  ThreadU
}

func NewUsecases(repo *repository.Repository) *Usecases {
	forumU := NewForumUsecase(repo.ForumRepository, repo.ThreadRepository)
	userU := NewUserUsecase(repo.UserRepository)
	threadU := NewThreadUseCase(repo.ThreadRepository)

	return &Usecases{
		ForumUsecase:   forumU,
		PostUsecase:    NewPostUsecase(repo.PostRepository, forumU, userU, threadU),
		UserUsecase:    userU,
		ServiceUsecase: NewServiceUsecase(repo.ServiceRepository),
		ThreadUsecase:  threadU,
	}
}

type ForumU interface {
	Create(forum *models.Forum) (createdForum *models.Forum, err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (createdThread *models.Thread, err error)
	GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error)
	GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error)
}

type PostU interface {
	Get(id int, details []string) (postDetailed *models.ParamsPost, err error)
	Update(post *models.Post) (updatedPost *models.Post, err error)
}

type ServiceU interface {
	Clear() (err error)
	Status() (status *models.ForumStatus, err error)
}

type ThreadU interface {
	Get(slugOrId string) (thread *models.Thread, err error)
	Update(slugOrId string, thread *models.Thread) (updatedThread *models.Thread, err error)
	Vote(slugOrId string, vote *models.Vote) (thread *models.Thread, err error)
	CreatePosts(slugOrId string, posts []*models.Post) (createdPosts []*models.Post, err error)
	GetPosts(slugOrId string, params *models.PostsQueryParams) (posts []*models.Post, err error)
}

type UserU interface {
	Get(nickname *string) (user *models.User, err error)
	All() (users *[]models.User, err error)
	Create(user *models.User) (users []*models.User, err error)
	Update(user *models.User) (updatedUser *models.User, err error)
}
