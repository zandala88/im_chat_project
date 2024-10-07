package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"im/model"
	"im/util"
)

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	// 参数验证
	friendIdStr := c.PostForm("friend_id")
	friendId := cast.ToInt64(friendIdStr)
	if friendId == 0 {
		zap.S().Error("AddFriend 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 获取自己的信息
	userId := util.GetUid(c)
	if userId == friendId {
		zap.S().Error("AddFriend 不能添加自己为好友")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询用户是否存在
	ub, err := model.GetUserById(friendId)
	if err != nil {
		zap.S().Error("AddFriend 好友不存在")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询是否已建立好友关系
	isFriend, err := model.IsFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("AddFriend 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if isFriend {
		zap.S().Error("AddFriend 请勿重复添加")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 建立好友关系
	err = model.CreateFriend(&model.Friend{
		UserID:   userId,
		FriendID: ub.ID,
	}, &model.Friend{
		UserID:   ub.ID,
		FriendID: userId,
	})
	if err != nil {
		zap.S().Error("AddFriend 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func FriendList(c *gin.Context) {
	userId := util.GetUid(c)
	friends, err := model.GetFriends(userId)
	if err != nil {
		zap.S().Error("FriendList 获取好友列表失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, gin.H{
		"friends": friends,
	})
}

func DeleteFriend(c *gin.Context) {
	// 参数验证
	friendIdStr := c.PostForm("friend_id")
	friendId := cast.ToInt64(friendIdStr)
	if friendId == 0 {
		zap.S().Error("DeleteFriend 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 获取自己的信息
	userId := util.GetUid(c)
	if userId == friendId {
		zap.S().Error("DeleteFriend 不能删除自己")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询用户是否存在
	ub, err := model.GetUserById(friendId)
	if err != nil {
		zap.S().Error("DeleteFriend 好友不存在")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询是否已建立好友关系
	isFriend, err := model.IsFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("DeleteFriend 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isFriend {
		zap.S().Error("DeleteFriend 请勿重复删除")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 删除好友关系
	err = model.DeleteFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("DeleteFriend 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}
