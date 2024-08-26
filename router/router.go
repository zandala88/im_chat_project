package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "im/docs"
	"im/middleware"
	"im/public"
	"im/service"
	"io"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()

	r := gin.Default()
	r.Use(middleware.Cors(), public.GinLogger(), public.GinRecovery())

	// Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Group("/")
	{
		r.POST("user/login", service.Login)
		r.POST("user/register", service.Register)
		// chat
		r.GET("/chat", service.Chat)
	}

	api := r.Group("/", middleware.Auth())
	{
		api.POST("add/friend", service.AddFriend)
		api.GET("get/friend", service.GetFriend)
		api.GET("get/friend/list", service.GetFriendList)
		//api.DELETE("friend", service.DeleteFriend)
	}
	return r
}
