package util

import "github.com/gin-gonic/gin"

func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func FailResp(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": -1,
		"msg":  msg,
	})
}
