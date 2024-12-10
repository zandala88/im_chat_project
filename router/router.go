package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"im/config"
	"im/public/middlewares"
	"im/service"
	"net/http"
)

// HTTPRouter http 路由
func HTTPRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 用户注册
	r.POST("/register", service.Register)

	// 用户登录
	r.POST("/login", service.Login)

	auth := r.Group("", middlewares.AuthCheck())
	{
		// 添加好友
		auth.POST("/friend", service.AddFriend)

		// 好友列表
		auth.GET("/friend/list", service.FriendList)

		// 删除好友
		auth.DELETE("/friend", service.DeleteFriend)

		// 群聊列表
		auth.GET("/group/list", service.GroupList)

		// 加入群聊
		auth.POST("/group/join", service.JoinGroup)

		// 退出群聊
		auth.POST("/group/exit", service.ExitGroup)

		// 解散群聊
		auth.DELETE("/group", service.DeleteGroup)

		// 创建群聊
		auth.POST("/group", service.CreateGroup)

		// 获取群成员列表
		auth.GET("/group/user/list", service.GroupUserList)
	}

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.HTTPServerPort)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
