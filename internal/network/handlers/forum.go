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

type HandlerForum struct {
	UseCase usecase.ForumU
}

func NewForumsHandler(useCase usecase.ForumU) *HandlerForum {
	return &HandlerForum{UseCase: useCase}
}

func (handler *HandlerForum) Get(c *gin.Context) {
	slug := c.Param("slug")
	if v, _ := check.GetInstance(); !v.CheckSlug(slug) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug"))
		return
	}

	forum, err := handler.UseCase.Get(slug)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
}

func (handler *HandlerForum) Create(c *gin.Context) {
	forum := &models.Forum{}

	err := easyjson.UnmarshalFromReader(c.Request.Body, forum)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	if v, _ := check.GetInstance(); !v.CheckSlug(forum.Slug) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug"))
		return
	}

	createdForum, err := handler.UseCase.Create(forum)

	if err != nil {
		if err.(errors.MsgErrors).Code() == errors.ForumAlreadyExists.Code() {
			c.JSON(errors.ForumAlreadyExists.Code(), createdForum)
			return
		}
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusCreated, createdForum)
}

func (handler *HandlerForum) GetUsers(c *gin.Context) {
	var slug string
	slug = c.Param("slug")
	if v, _ := check.GetInstance(); !v.CheckSlug(slug) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug"))
		return
	}

	params := &models.ForumUserQueryParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректные query params"))
	}

	v, _ := check.GetInstance()
	v.CheckForumUserQuery(params)

	threads, err := handler.UseCase.GetUsers(slug, params)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, threads)
	return
}

func (handler *HandlerForum) CreateThread(c *gin.Context) {
	thread := &models.Thread{}

	thread.Forum = c.Param("slug")
	if v, _ := check.GetInstance(); !v.CheckSlug(thread.Forum) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug forum"))
		return
	}

	err := easyjson.UnmarshalFromReader(c.Request.Body, thread)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	if v, _ := check.GetInstance(); thread.Slug != "" && !v.CheckSlug(thread.Slug) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug thread"))
		return
	}

	createdThread, err := handler.UseCase.CreateThread(thread)
	if err != nil {
		if err.(errors.MsgErrors).Code() == errors.ThreadAlreadyExists.Code() {
			c.JSON(errors.ThreadAlreadyExists.Code(), createdThread)
			return
		}
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusCreated, createdThread)
	return
}

func (handler *HandlerForum) GetThreads(c *gin.Context) {
	slug := c.Param("slug")
	if v, _ := check.GetInstance(); !v.CheckSlug(slug) {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректный slug"))
		return
	}

	params := &models.ForumQueryParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest.SetTextDetails("Не корректные query params"))
	}

	v, _ := check.GetInstance()
	v.CheckForumQuery(params)

	threads, err := handler.UseCase.GetThreads(slug, params)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, threads)
	return
}
