package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "im/api/v1"
	"im/util"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()
	server.Use(Cors())
	server.Use(Recovery)

	socket := RunSocekt

	group := server.Group("")
	{
		group.GET("/user", v1.GetUserList)
		group.GET("/user/:uuid", v1.GetUserDetails)
		group.GET("/user/name", v1.GetUserOrGroupByName)
		group.POST("/user/register", v1.Register)
		group.POST("/user/login", v1.Login)
		group.PUT("/user", v1.ModifyUserInfo)

		group.POST("/friend", v1.AddFriend)

		group.GET("/message", v1.GetMessage)

		group.GET("/file/:fileName", v1.GetFile)
		group.POST("/file", v1.SaveFile)

		group.GET("/group/:uuid", v1.GetGroup)
		group.POST("/group/:uuid", v1.SaveGroup)
		group.POST("/group/join/:userUuid/:groupUuid", v1.JoinGroup)
		group.GET("/group/user/:uuid", v1.GetGroupUsers)

		group.GET("/socket.io", socket)
	}
	return server
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			util.SuccessResp(c, nil)
		}

		defer func() {
			if err := recover(); err != nil {
				zap.S().Error("Http err = ", err)
			}
		}()

		c.Next()
	}
}

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Error("gin catch err =  ", r)
			util.FailResp(c, "系统内部错误")
		}
	}()
	c.Next()
}
