package controller

import (
	"github.com/gin-gonic/gin"
	"im/model"
	"im/service"
	"im/util"
	"strconv"
)

// 创建分组参数
type CommunityParam struct {
	Name    string `json:"name"`
	OwnerId int    `json:"ownerId"`
}

// 获取用户所加的群
func Communities(ctx *gin.Context) {
	auth, _ := ctx.Get("auth")

	user := auth.(model.User)

	communities, err := service.UserCommunities(user.ID)
	if err != nil {
		h := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		h.Fail(ctx.Writer)
	}

	h := util.Response{
		Code: 0,
		Data: communities,
	}

	h.Success(ctx.Writer)
}

// 获取群组信息
func CommunityInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	groupId, _ := strconv.Atoi(id)

	community, err := service.CommunityInfo(&model.Community{
		ID: groupId,
	})

	if err != nil {
		h := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		h.Fail(ctx.Writer)
		return
	}

	h := util.Response{
		Code: 0,
		Data: community,
	}

	h.Success(ctx.Writer)
}

// 创建群组功能
func CreateCommunity(ctx *gin.Context) {
	param := CommunityParam{}
	ctx.ShouldBindJSON(&param)
	community, err := service.CreateCommunity(param.OwnerId, param.Name)
	if err != nil {
		h := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		h.Fail(ctx.Writer)
		return
	}

	h := util.Response{
		Code: 0,
		Data: community,
	}

	h.Success(ctx.Writer)
}

// 加入群聊
func JoinCommunity(ctx *gin.Context) {
	groupId, _ := ctx.Params.Get("id")
	id, _ := strconv.Atoi(groupId)
	auth, _ := ctx.Get("auth")
	user := auth.(model.User)
	ok, err := service.JoinCommunity(id, user.ID)
	if err != nil || !ok {
		h := util.Response{
			Code:    -1,
			Message: err.Error(),
		}
		h.Fail(ctx.Writer)
		return
	}

	h := util.Response{
		Code:    0,
		Message: "加入成功",
	}
	h.Success(ctx.Writer)
}

// 获取群组的信息(最新的 10 条)
func GetCommunityMessages(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, _ := strconv.Atoi(communityId)
	messages := service.GetCommunityMessages(id)

	h := util.Response{
		Code: 0,
		Data: messages,
	}

	h.Success(ctx.Writer)
}
