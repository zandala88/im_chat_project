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
		auth.POST("/friend/add", service.AddFriend)

		// 创建群聊
		auth.POST("/group/create", service.CreateGroup)

		// 获取群成员列表
		auth.GET("/group_user/list", service.GroupUserList)
	}

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.HTTPServerPort)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
