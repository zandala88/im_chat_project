package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "im/docs"
	"im/middleware"
	"im/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	gin.DisableConsoleColor()
	r.Use(middleware.Cors())
	// Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Group("/")
	{
		r.POST("user/login", service.Login)
		r.POST("user/register", service.Register)
	}

	api := r.Group("/", middleware.Auth())
	{
		api.POST("add/friend", service.AddFriend)
		api.GET("get/friend", service.GetFriend)
		api.GET("get/friend/list", service.GetFriendList)
		//api.DELETE("friend", service.DeleteFriend)

		// chat
		api.GET("/chat", service.Chat)
	}
	return r
}
