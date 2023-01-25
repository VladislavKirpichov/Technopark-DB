package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v.kirpichov/db_tp/internal/usecase"
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
)

type HandlerServices struct {
	UseCase usecase.ServiceU
}

func NewServicesHandler(useCase usecase.ServiceU) *HandlerServices {
	return &HandlerServices{UseCase: useCase}
}

func (handler *HandlerServices) Clear(c *gin.Context) {
	err := handler.UseCase.Clear()
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (handler *HandlerServices) Status(c *gin.Context) {
	status, err := handler.UseCase.Status()
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, status)
}
