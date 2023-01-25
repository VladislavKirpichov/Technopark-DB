package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/usecase"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type HandlerPosts struct {
	UseCase usecase.PostU
}

func NewPostsHandler(useCase usecase.PostU) *HandlerPosts {
	return &HandlerPosts{UseCase: useCase}
}

func (handler *HandlerPosts) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	detailsRaw := c.Query("related")
	var details []string
	if detailsRaw != "" {
		details = strings.Split(detailsRaw, ",")
	}

	post, err := handler.UseCase.Get(int(id), details)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, post)
	return
}

func (handler *HandlerPosts) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	post := &models.Post{}
	err = easyjson.UnmarshalFromReader(c.Request.Body, post)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	post.ID = int(id)

	forum, err := handler.UseCase.Update(post)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
	return
}
