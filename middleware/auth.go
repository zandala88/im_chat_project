package middleware

import (
	"github.com/gin-gonic/gin"
	"im/model"
	"im/service"
	"im/util"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" || len(token) == 0 {
			h := util.Response{
				Code:    403,
				Message: "用户未登录",
			}
			h.Fail(ctx.Writer)
			return
		}
		// 获取用户
		user := model.User{}
		service.DbEngine.Where("token = ?", token).First(&user)
		if user.ID == 0 {
			h := util.Response{
				Code:    403,
				Message: "无效的 token",
			}
			h.Fail(ctx.Writer)
			return
		}
		ctx.Set("auth", user)
		ctx.Next()
	}
}
