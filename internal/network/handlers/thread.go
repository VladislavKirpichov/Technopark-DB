package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/usecase"
	"github.com/v.kirpichov/db_tp/internal/utils/check"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type HandlerThreads struct {
	UseCase usecase.ThreadU
}

func NewThreadsHandler(useCase usecase.ThreadU) *HandlerThreads {
	return &HandlerThreads{UseCase: useCase}
}

func (handler *HandlerThreads) Get(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	forum, err := handler.UseCase.Get(slugOrId)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
	return
}

func (handler *HandlerThreads) Update(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	thread := &models.Thread{}
	err := easyjson.UnmarshalFromReader(c.Request.Body, thread)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	forum, err := handler.UseCase.Update(slugOrId, thread)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
	return
}

func (handler *HandlerThreads) Vote(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	vote := &models.Vote{}
	err := easyjson.UnmarshalFromReader(c.Request.Body, vote)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	forum, err := handler.UseCase.Vote(slugOrId, vote)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
	return
}

func (handler *HandlerThreads) PostsCreate(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	var posts models.Posts
	err := easyjson.UnmarshalFromReader(c.Request.Body, &posts)

	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	createdPosts, err := handler.UseCase.CreatePosts(slugOrId, posts)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusCreated, createdPosts)
	return
}

func (handler *HandlerThreads) GetPosts(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	params := &models.PostsQueryParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректные query params"))
	}

	if v, _ := check.GetInstance(); !v.CheckPostsQuery(params) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректные query params"))
	}

	createdPosts, err := handler.UseCase.GetPosts(slugOrId, params)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, createdPosts)
	return
}
