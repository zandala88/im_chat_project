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
		zap.S().Error("[CreateGroup] 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}
	zap.S().Debug("[CreateGroup] [idsStr] = ", idsStr)
	ids := make([]int64, 0, len(idsStr)+1)
	for i := range idsStr {
		ids = append(ids, cast.ToInt64(idsStr[i]))
	}

	// 获取用户信息
	userId := util.GetUid(c)
	ids = append(ids, userId)

	// 获取 ids 用户信息
	userRepo := model.NewUserRepo(c)
	ids, err := userRepo.GetUserIdByIds(ids)
	if err != nil {
		zap.S().Error("[CreateGroup] [model.GetUserIdByIds] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	zap.S().Info("[CreateGroup] [ids] = ", ids)

	// 创建群组
	group := &model.Group{
		Name:    name,
		OwnerID: userId,
	}
	groupRepo := model.NewGroupRepo(c)
	err = groupRepo.CreateGroup(group, ids)
	if err != nil {
		zap.S().Error("[CreateGroup] [model.CreateGroup] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 将群成员信息更新到 Redis
	err = cache.SetGroupUser(group.ID, ids)
	if err != nil {
		zap.S().Error("[CreateGroup] [model.SetGroupUser] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, gin.H{
		"id": cast.ToString(group.ID),
	})
}

func GroupList(c *gin.Context) {
	userId := util.GetUid(c)
	groupRepo := model.NewGroupRepo(c)
	groups, err := groupRepo.GetGroups(userId)
	if err != nil {
		zap.S().Error("[GroupList] [model.GetGroups] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, groups)
}

func JoinGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("[JoinGroup] groupId == 0")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	groupRepo := model.NewGroupRepo(c)
	_, err := groupRepo.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("[JoinGroup] [model.GetGroupById] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否已经在群里
	userId := util.GetUid(c)
	groupUserRepo := model.NewGroupUserRepo(c)
	isBelong, err := groupUserRepo.IsBelongToGroup(userId, groupId)
	if err != nil {
		zap.S().Error("[JoinGroup] [model.IsBelongToGroup] 获取群组失败 [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if isBelong {
		zap.S().Error("[JoinGroup] isBelong == true")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 加入群组
	err = groupUserRepo.JoinGroup(groupId, userId)
	if err != nil {
		zap.S().Error("[JoinGroup] [model.JoinGroup] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 将群成员信息更新到 Redis
	err = cache.SetGroupUser(groupId, []int64{userId})
	if err != nil {
		zap.S().Error("[JoinGroup] [cache.SetGroupUser] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func ExitGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("[ExitGroup] groupId == 0")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	groupRepo := model.NewGroupRepo(c)
	_, err := groupRepo.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("[ExitGroup] [model.GetGroupById] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否在群里
	userId := util.GetUid(c)
	groupUserRepo := model.NewGroupUserRepo(c)
	isBelong, err := groupUserRepo.IsBelongToGroup(userId, groupId)
	if err != nil {
		zap.S().Error("[ExitGroup] [model.IsBelongToGroup] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isBelong {
		zap.S().Error("[ExitGroup] isBelong == false")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 退出群组
	err = groupUserRepo.ExitGroup(groupId, userId)
	if err != nil {
		zap.S().Error("[ExitGroup] [model.ExitGroup] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 从 Redis 中删除群成员信息
	err = cache.DeleteGroupUser(groupId, userId)
	if err != nil {
		zap.S().Error("[ExitGroup] [cache.DeleteGroupUser] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}

func DeleteGroup(c *gin.Context) {
	groupIdStr := c.PostForm("group_id")
	groupId := cast.ToInt64(groupIdStr)
	if groupId == 0 {
		zap.S().Error("[DeleteGroup] groupId == 0")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查看群是否存在
	groupRepo := model.NewGroupRepo(c)
	_, err := groupRepo.GetGroupById(groupId)
	if err != nil {
		zap.S().Error("[DeleteGroup] [model.GetGroupById] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 查看用户是否是群主
	userId := util.GetUid(c)
	isOwner, err := groupRepo.IsGroupOwner(userId, groupId)
	if err != nil {
		zap.S().Error("[DeleteGroup] [model.IsGroupOwner] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if !isOwner {
		zap.S().Error("[DeleteGroup] isOwner == false")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 删除群组
	err = groupRepo.DeleteGroup(groupId)
	if err != nil {
		zap.S().Error("[DeleteGroup] [model.DeleteGroup] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 从 Redis 中删除群成员信息
	err = cache.DeleteGroupUserAll(groupId)
	if err != nil {
		zap.S().Error("[DeleteGroup] [cache.DeleteGroupUserAll] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, nil)
}
