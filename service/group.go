package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"im/model"
	"im/server/cache"
	"im/util"
)

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	// 参数校验
	name := c.PostForm("name")
	idsStr := c.PostFormArray("ids") // 群成员 id，不包括群创建者
	if name == "" || len(idsStr) == 0 {
		zap.S().Error("CreateGroup 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}
	ids := make([]int64, 0, len(idsStr)+1)
	for i := range idsStr {
		ids = append(ids, cast.ToInt64(idsStr[i]))
	}

	// 获取用户信息
	userId := util.GetUid(c)
	ids = append(ids, userId)

	// 获取 ids 用户信息
	ids, err := model.GetUserIdByIds(ids)
	if err != nil {
		zap.S().Error("CreateGroup 获取用户信息失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 创建群组
	group := &model.Group{
		Name:    name,
		OwnerID: userId,
	}
	err = model.CreateGroup(group, ids)
	if err != nil {
		zap.S().Error("CreateGroup 创建群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 将群成员信息更新到 Redis
	err = cache.SetGroupUser(group.ID, ids)
	if err != nil {
		zap.S().Error("CreateGroup 缓存群成员失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, gin.H{
		"id": cast.ToString(group.ID),
	})
}

func GroupList(c *gin.Context) {
	userId := util.GetUid(c)
	groups, err := model.GetGroups(userId)
	if err != nil {
		zap.S().Error("GroupList 获取群组列表失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, groups)
}

func JoinGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("JoinGroup 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	_, err := model.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("JoinGroup 获取群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否已经在群里
	userId := util.GetUid(c)
	isBelong, err := model.IsBelongToGroup(userId, groupId)
	if err != nil {
		zap.S().Error("JoinGroup 获取群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if isBelong {
		zap.S().Error("JoinGroup 用户已在群组中")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 加入群组
	err = model.JoinGroup(groupId, userId)
	if err != nil {
		zap.S().Error("JoinGroup 加入群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 将群成员信息更新到 Redis
	err = cache.SetGroupUser(groupId, []int64{userId})
	if err != nil {
		zap.S().Error("JoinGroup 缓存群成员失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func ExitGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("ExitGroup 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	_, err := model.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("ExitGroup 获取群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否在群里
	userId := util.GetUid(c)
	isBelong, err := model.IsBelongToGroup(userId, groupId)
	if err != nil {
		zap.S().Error("ExitGroup 获取群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isBelong {
		zap.S().Error("ExitGroup 用户不在群组中")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 退出群组
	err = model.ExitGroup(groupId, userId)
	if err != nil {
		zap.S().Error("ExitGroup 退出群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 从 Redis 中删除群成员信息
	err = cache.DeleteGroupUser(groupId, userId)
	if err != nil {
		zap.S().Error("ExitGroup 删除群成员失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func DeleteGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("DeleteGroup 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	_, err := model.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("DeleteGroup 获取群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否是群主
	userId := util.GetUid(c)
	isOwner, err := model.IsGroupOwner(userId, groupId)
	if err != nil {
		zap.S().Error("DeleteGroup 获取群主失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isOwner {
		zap.S().Error("DeleteGroup 用户不是群主")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 删除群组
	err = model.DeleteGroup(groupId)
	if err != nil {
		zap.S().Error("DeleteGroup 删除群组失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 从 Redis 中删除群成员信息
	err = cache.DeleteGroupUserAll(groupId)
	if err != nil {
		zap.S().Error("DeleteGroup 删除群成员失败", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}
