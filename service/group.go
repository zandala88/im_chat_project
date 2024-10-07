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
