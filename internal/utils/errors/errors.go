package errors

import (
	"net/http"

	"github.com/v.kirpichov/db_tp/internal/models"
)

type MsgErrors interface {
	Error() string
	SetTextDetails(text string) *models.Message
	Code() int
}

var (
	Err23503 = "23503"
	Err23502 = "23502"
	Err23505 = "23505"
	P0001    = "P0001"
)

var (
	ServerInternal MsgErrors = &models.Message{ErrorCode: http.StatusInternalServerError, Msg: "internal server error"}
	BadRequest     MsgErrors = &models.Message{ErrorCode: http.StatusBadRequest, Msg: "bad request"}
)

var (
	ConflictUserCreate MsgErrors = &models.Message{ErrorCode: http.StatusConflict, Msg: "пользователь c таким email или nickname уже существует"}
	NotFoundUser       MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "не найден юзер"}
	ConflictUserUpdate MsgErrors = &models.Message{ErrorCode: http.StatusConflict, Msg: "новые данные профиля пользователя конфликтуют с имеющимися пользователями"}
	NotFoundUserUpdate MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "не найден пользователь для обновления"}
)

var (
	ThreadAlreadyExists        MsgErrors = &models.Message{ErrorCode: http.StatusConflict, Msg: "тред уже присутсвует в базе данных"}
	ThreadUpdateNotFound       MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "не найден тред для обновления"}
	ThreadUserOrThreadNotFound MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "не найден пользователь или тред для голосования"}
	ThreadUserOrForumNotFound  MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "автор треда или форуи не найдены"}
	ThreadNotFound             MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "тред не найден"}
)

var (
	NotFoundForumUser  MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "владелец форума не найден"}
	ForumAlreadyExists MsgErrors = &models.Message{ErrorCode: http.StatusConflict, Msg: "форум уже присутсвует в базе данных"}
	NotFoundForum      MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "форум не найден"}
)

var (
	PostWrongParent  MsgErrors = &models.Message{ErrorCode: http.StatusConflict, Msg: "не найден указанный родетель в данном треде"}
	PostUserNotFound MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "автор поста не найден"}
	PostNotFound     MsgErrors = &models.Message{ErrorCode: http.StatusNotFound, Msg: "не найден пост для обновления"}
)
