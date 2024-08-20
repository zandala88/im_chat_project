package controller

import (
	"github.com/gin-gonic/gin"
	"im/model"
	"im/service"
	"im/util"
	"net/http"
	"strconv"
)

// 添加好友方法
func AddFriend(ctx *gin.Context) {
	friendId, _ := strconv.Atoi(ctx.PostForm("friend_id"))
	userId, _ := strconv.Atoi(ctx.PostForm("user_id"))

	contact, err := service.AddFriend(int(userId), int(friendId))
	if err != nil {
		r := util.Response{
			Code:    500,
			Message: err.Error(),
		}

		r.Fail(ctx.Writer)
	} else {
		r := util.Response{
			Code:    0,
			Message: "添加好友成功",
			Data:    contact,
		}
		r.Success(ctx.Writer)
	}
}

// 删除好友
func DeleteFriend(ctx *gin.Context) {
	friendId, _ := strconv.Atoi(ctx.PostForm("friend_id"))
	userId, _ := strconv.Atoi(ctx.PostForm("user_id"))

	isOk, err := service.DeleteFriend(int(userId), int(friendId))
	if !isOk || err != nil {
		r := util.Response{
			Code:    -1,
			Message: err.Error(),
		}

		r.Fail(ctx.Writer)
	}

	r := util.Response{
		Code:    0,
		Message: "操作成功",
	}

	r.Success(ctx.Writer)
}

// 获取用户的要有列表
func GetFriends(ctx *gin.Context) {
	auth, ok := ctx.Get("auth")
	if !ok {
		h := util.Response{
			Code:    http.StatusForbidden,
			Message: "用户未登录",
		}
		h.Fail(ctx.Writer)
	}
	// 查找用户的好友
	user := auth.(model.User)
	users, err := service.Friends(user)
	if err != nil {
		r := util.Response{
			Code:    1,
			Message: err.Error(),
		}
		r.Fail(ctx.Writer)
	}

	r := util.Response{
		Code: 0,
		Data: users,
	}

	r.Success(ctx.Writer)
}

// 获取好友信息
func GetFriend(ctx *gin.Context) {
	friendId, ok := ctx.Params.Get("friendId")
	if !ok {
		r := util.Response{
			Code:    -1,
			Message: "参数无效",
		}
		r.Fail(ctx.Writer)
	}
	fId, err := strconv.Atoi(friendId)
	if err != nil {
		r := util.Response{
			Code:    -1,
			Message: "参数无效",
		}
		r.Fail(ctx.Writer)
	}
	friend, err := service.GetFriend(fId)
	if err != nil {
		r := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		r.Fail(ctx.Writer)
	}

	r := util.Response{
		Code: 0,
		Data: friend,
	}
	r.Success(ctx.Writer)
}
