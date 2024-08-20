package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 跨域中间件
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "DNT,Keep-Alive,User-Agent,Cache-Control,Content-Type,Authorization,x-token")
		// 设置 header 为 json, 返回 json
		context.Writer.Header().Set("Content-Encoding", "zh-CN")
		context.Writer.Header().Set("Content-Type", "application/json")
		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "ok")
			return
		}
		context.Next()
	}
}
