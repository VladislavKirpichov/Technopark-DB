package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository/impl"
)

type Repository struct {
	ForumRepository   ForumR
	PostRepository    PostR
	UserRepository    UserR
	ServiceRepository ServiceR
	ThreadRepository  ThreadR
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		ForumRepository:   impl.NewForumRepository(db),
		PostRepository:    impl.NewPostRepository(db),
		UserRepository:    impl.NewUserRepository(db),
		ServiceRepository: impl.NewServiceRepository(db),
		ThreadRepository:  impl.NewThreadRepository(db),
	}
}

type ForumR interface {
	Create(forum *models.Forum) (createdForum *models.Forum, err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (createdThread *models.Thread, err error)
	GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error)
	GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error)
}

type PostR interface {
	Get(id int) (post *models.Post, err error)
	Update(post *models.Post) (updatedPost *models.Post, err error)
}

type UserR interface {
	Get(nickname *string) (user *models.User, err error)
	Update(user *models.User) (updatedUser *models.User, err error)
	GetUsersByUserNicknameOrEmail(user *models.User) (users []*models.User, err error)
	All() (users *[]models.User, err error)
	Create(user *models.User) (err error)
}

type ServiceR interface {
	Clear() (err error)
	Status() (status *models.ForumStatus, err error)
}

type ThreadR interface {
	GetBySlug(slug string) (*models.Thread, error)
	UpdateByID(thread *models.Thread) (*models.Thread, error)
	GetByID(id int) (*models.Thread, error)
	CreatePosts(threadId int, forumSlug string, post []*models.Post) (createdPosts []*models.Post, err error)
	GetPosts(threadId int, params *models.PostsQueryParams) (posts []*models.Post, err error)
	VoteBySlug(slug string, vote *models.Vote) (err error)
	UpdateBySlug(thread *models.Thread) (*models.Thread, error)
	VoteByID(threadId int, vote *models.Vote) (err error)
}
