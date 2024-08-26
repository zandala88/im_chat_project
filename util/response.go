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

func FailRespWithCode(c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  errMsg[code],
	})
}

var errMsg = map[int]string{
	NoToken:             "没有token,请重新登录",
	TokenError:          "token错误",
	InvalidToken:        "无效的Token",
	ShouldBindJSONError: "参数绑定错误",
	CURDSelectError:     "查询失败",
	CURDInsertError:     "插入失败",
	CURDUpdateError:     "更新失败",
	CURDDeleteError:     "删除失败",
	SendEmailIn3Min:     "3分钟内只能发送一次邮件",
	SendEmailError:      "发送邮件失败",
}

const (
	NoToken             = 2003
	TokenError          = 2004
	InvalidToken        = 2005
	ShouldBindJSONError = 2006
	CURDSelectError     = 2007
	CURDInsertError     = 2008
	CURDUpdateError     = 2009
	CURDDeleteError     = 2010
	SendEmailIn3Min     = 2011
	SendEmailError      = 2012
)
