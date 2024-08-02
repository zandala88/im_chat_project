package router

import (
	"github.com/gin-gonic/gin"
	"im_chat_project/api/handler"
	"im_chat_project/api/middleware"
	"im_chat_project/public"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	initUserRouter(r)
	initPushRouter(r)
	r.NoRoute(func(c *gin.Context) {
		public.ResponseError(c, 40001, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/login", handler.Login)
	userGroup.POST("/register", handler.Register)
	userGroup.Use(middleware.CheckSessionIdMiddleware())
	{
		userGroup.POST("/checkAuth", handler.CheckAuth)
		userGroup.POST("/logout", handler.Logout)
	}

}

func initPushRouter(r *gin.Engine) {
	pushGroup := r.Group("/push")
	pushGroup.Use(middleware.CheckSessionIdMiddleware())
	{
		pushGroup.POST("/push", handler.Push)
		pushGroup.POST("/pushRoom", handler.PushRoom)
		pushGroup.POST("/count", handler.Count)
		pushGroup.POST("/getRoomInfo", handler.GetRoomInfo)
	}

}
