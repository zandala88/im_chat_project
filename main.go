package main

import (
	"github.com/gin-gonic/gin"
	"im/controller"
	"im/middleware"
)

// Router of gin
var Router *gin.Engine

func init() {
	gin.DisableConsoleColor()
	Router = gin.Default()

	Router.Use(middleware.Cors())

	Router.Group("/")
	{
		Router.POST("user/login", controller.Login)
		Router.POST("user/register", controller.Register)
	}

	api := Router.Group("/", middleware.Auth())
	{
		api.POST("friend", controller.AddFriend)
		api.GET("friend/:friendId", controller.GetFriend)
		api.GET("friends", controller.GetFriends)
		api.DELETE("friend", controller.DeleteFriend)

		// chat
		api.GET("/chat", controller.Chat)

		// 群组功能
		api.GET("/community/:id", controller.CommunityInfo)
		api.POST("/community", controller.CreateCommunity)
		// 获取用户添加的群组功能
		api.GET("/communities", controller.Communities)
		api.POST("/join/community/:id", controller.JoinCommunity)

		// 获取好友的聊天信息
		api.GET("messages/:friendId", controller.GetFriendMessages)
		api.GET("community_messages/:communityId", controller.GetCommunityMessages)
	}
}
func main() {
	//// 绑定请求和处理函数
	//http.HandleFunc("/user/login", controller.Login)
	//http.HandleFunc("/user/register", controller.Register)
	//http.HandleFunc("/add/friend", controller.AddFriend)
	//http.HandleFunc("/delete/friend", controller.DeleteFriend)

	// 提供静态资源目录的支持
	//http.Handle("/", http.FileServer(http.Dir("./"))) // 此方法会暴露当前项目的所有文件, 所以需要制定暴露的目录
	//http.Handle("/assets/", http.FileServer(http.Dir(".")))

	// 启动 web 服务时期
	//http.ListenAndServe(":8080", nil)

	Router.Run(":8080")
}
