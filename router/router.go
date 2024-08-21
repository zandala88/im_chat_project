package router

import (
	"github.com/gin-gonic/gin"
	"im/middleware"
	"im/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	gin.DisableConsoleColor()

	r.Use(middleware.Cors())

	r.Group("/")
	{
		r.POST("user/login", service.Login)
		r.POST("user/register", service.Register)
	}

	//api := r.Group("/", middleware.Auth())
	//{
	//	api.POST("friend", controller.AddFriend)
	//	api.GET("friend/:friendId", controller.GetFriend)
	//	api.GET("friends", controller.GetFriends)
	//	api.DELETE("friend", controller.DeleteFriend)
	//
	//	// chat
	//	api.GET("/chat", controller.Chat)
	//
	//	// 群组功能
	//	api.GET("/community/:id", controller.CommunityInfo)
	//	api.POST("/community", controller.CreateCommunity)
	//	// 获取用户添加的群组功能
	//	api.GET("/communities", controller.Communities)
	//	api.POST("/join/community/:id", controller.JoinCommunity)
	//
	//	// 获取好友的聊天信息
	//	api.GET("messages/:friendId", controller.GetFriendMessages)
	//	api.GET("community_messages/:communityId", controller.GetCommunityMessages)
	//}
	return r
}
