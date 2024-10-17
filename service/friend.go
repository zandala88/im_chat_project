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
		zap.S().Error("[AddFriend] friendId == 0 ")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 获取自己的信息
	userId := util.GetUid(c)
	if userId == friendId {
		zap.S().Error("[AddFriend] userId == friendId")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询用户是否存在
	ub, err := model.GetUserById(friendId)
	if err != nil {
		zap.S().Error("[AddFriend] [model.GetUserById] [err] = ", err.Error())
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询是否已建立好友关系
	isFriend, err := model.IsFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("[AddFriend] [model.IsFriend] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if isFriend {
		zap.S().Error("[AddFriend] isFriend == true")
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
		zap.S().Error("[AddFriend] [model.CreateFriend] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func FriendList(c *gin.Context) {
	userId := util.GetUid(c)
	friends, err := model.GetFriends(userId)
	if err != nil {
		zap.S().Error("[FriendList] [model.GetFriends] [err] = ", err.Error())
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
		zap.S().Error("[DeleteFriend] friend_id == 0")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 获取自己的信息
	userId := util.GetUid(c)
	if userId == friendId {
		zap.S().Error("[DeleteFriend] userId == friendId")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询用户是否存在
	ub, err := model.GetUserById(friendId)
	if err != nil {
		zap.S().Error("[DeleteFriend] [model.GetUserById] [err] = ", err.Error())
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询是否已建立好友关系
	isFriend, err := model.IsFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("[DeleteFriend] [model.IsFriend] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isFriend {
		zap.S().Error("[DeleteFriend] isFriend == false")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 删除好友关系
	err = model.DeleteFriend(userId, ub.ID)
	if err != nil {
		zap.S().Error("[DeleteFriend] [model.DeleteFriend] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}
