package middleware

import (
	"github.com/gin-gonic/gin"
	"im/util"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			util.FailRespWithCode(c, util.NoToken)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.Split(authHeader, ".")
		if len(parts) != 3 {
			util.FailRespWithCode(c, util.TokenError)
			c.Abort()
			return
		}
		mc, err := util.VerifyJWT(authHeader)
		if err != nil {
			util.FailRespWithCode(c, util.InvalidToken)
			c.Abort()
			return
		}

		c.Set("id", mc.Id)
		c.Next()
	}
}
