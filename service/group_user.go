package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"im/model"
	"im/server/cache"
	"im/util"
)

// GroupUserList 获取群成员列表
func GroupUserList(c *gin.Context) {
	// 参数校验
	groupIdStr := c.Query("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("GroupUserList 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}
	userId := util.GetUid(c)

	// 验证用户是否属于该群
	isBelong, err := model.IsBelongToGroup(userId, groupId)
	if err != nil {
		zap.S().Error("GroupUserList 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isBelong {
		zap.S().Error("GroupUserList 用户不属于该群")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 获取群成员id列表
	ids, err := GetGroupUser(groupId)
	if err != nil {
		zap.S().Error("GroupUserList 系统错误", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, cast.ToString(id))
	}
	util.SuccessResp(c, gin.H{
		"ids": idsStr,
	})
}

// GetGroupUser 获取群成员
// 从缓存中获取，如果缓存中没有，获取后加入缓存
func GetGroupUser(groupId int64) ([]int64, error) {
	userIds, err := cache.GetGroupUser(groupId)
	if err != nil {
		zap.S().Error("GetGroupUser 缓存错误", err.Error())
		return nil, err
	}
	if len(userIds) != 0 {
		return userIds, nil
	}

	userIds, err = model.GetGroupUserIdsByGroupId(groupId)
	if err != nil {
		zap.S().Error("GetGroupUser 获取群成员错误", err.Error())
		return nil, err
	}
	err = cache.SetGroupUser(groupId, userIds)
	if err != nil {
		zap.S().Error("GetGroupUser 缓存错误", err.Error())
		return nil, err
	}

	return userIds, nil
}
