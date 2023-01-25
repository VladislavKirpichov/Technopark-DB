package network

import (
	"github.com/gin-gonic/gin"
	"github.com/v.kirpichov/db_tp/internal/network/handlers"
)

var Urls = map[string]string{
	"Root":    "/api",
	"User":    "/user",
	"Forum":   "/forum",
	"Thread":  "/thread",
	"Service": "/service",
	"Post":    "/post",
}

func InitRoutes(handlers *handlers.Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	apiGroup := r.Group(Urls["Root"])

	userRouter := apiGroup.Group(Urls["User"])
	userRouter.GET("/:nickname/profile", handlers.UserHandler.Get)
	userRouter.POST("/:nickname/profile", handlers.UserHandler.Update)
	userRouter.POST("/:nickname/create", handlers.UserHandler.Create)

	forumRouter := apiGroup.Group(Urls["Forum"])
	forumRouter.GET("/:slug/details", handlers.ForumHandler.Get)
	forumRouter.POST("/create", handlers.ForumHandler.Create)
	forumRouter.GET("/:slug/users", handlers.ForumHandler.GetUsers)
	forumRouter.GET("/:slug/threads", handlers.ForumHandler.GetThreads)
	forumRouter.POST("/:slug/create", handlers.ForumHandler.CreateThread)

	threadRouter := apiGroup.Group(Urls["Thread"])
	threadRouter.GET("/:slug_or_id/details", handlers.ThreadHandler.Get)
	threadRouter.POST("/:slug_or_id/details", handlers.ThreadHandler.Update)
	threadRouter.POST("/:slug_or_id/vote", handlers.ThreadHandler.Vote)
	threadRouter.POST("/:slug_or_id/create", handlers.ThreadHandler.PostsCreate)
	threadRouter.GET("/:slug_or_id/posts", handlers.ThreadHandler.GetPosts)

	serviceRouter := apiGroup.Group(Urls["Service"])
	serviceRouter.POST("/clear", handlers.ServiceHandler.Clear)
	serviceRouter.GET("/status", handlers.ServiceHandler.Status)

	postRouter := apiGroup.Group(Urls["Post"])
	postRouter.GET("/:id/details", handlers.PostHandler.Get)
	postRouter.POST("/:id/details", handlers.PostHandler.Update)

	return r
}
