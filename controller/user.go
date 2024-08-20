package controller

import (
	"github.com/gin-gonic/gin"
	"im/service"
	"im/util"
)

type Params struct {
	Mobile   string `json:"mobile"`
	Passwd   string `json:"passwd"`
	Nickname string `json:"nickname"`
}

// 用户登录
func Login(ctx *gin.Context) {
	params := Params{}
	ctx.ShouldBindJSON(&params)
	user, err := service.UserService{}.Login(params.Mobile, params.Passwd)

	if err != nil {
		// 登录失败 json
		resp := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		resp.Fail(ctx.Writer)
	} else {
		resp := util.Response{
			Code: 1,
			Data: user,
		}

		resp.Success(ctx.Writer)
	}
}

func Register(ctx *gin.Context) {
	// 解析参数
	params := Params{}
	ctx.ShouldBindJSON(&params)

	user, err := service.UserService{}.Register(params.Mobile, params.Passwd, params.Nickname)
	if err != nil {
		resp := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		resp.Fail(ctx.Writer)
	} else {
		resp := util.Response{
			Code:    1,
			Message: "注册成功",
			Data:    user,
		}

		resp.Success(ctx.Writer)
	}
}
