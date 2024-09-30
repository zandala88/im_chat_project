package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"im/internal/model"
	"im/internal/service"
	"im/util"
)

// GetGroup
// 获取分组列表
func GetGroup(c *gin.Context) {
	uuid := c.Param("uuid")
	groups, err := service.GroupService.GetGroups(uuid)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, groups)
}

// SaveGroup
// 保存分组列表
func SaveGroup(c *gin.Context) {
	uuid := c.Param("uuid")

	group := &model.Group{}
	err := c.ShouldBindJSON(&group)
	if err != nil {
		zap.S().Error("ShouldBindJSON err = ", err)
		util.FailResp(c, err.Error())
		return
	}

	err = service.GroupService.SaveGroup(uuid, group)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, nil)
}

// JoinGroup
// 加入组别
func JoinGroup(c *gin.Context) {
	userUuid := c.Param("userUuid")
	groupUuid := c.Param("groupUuid")
	err := service.GroupService.JoinGroup(groupUuid, userUuid)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, nil)
}

// GetGroupUsers
// 获取组内成员信息
func GetGroupUsers(c *gin.Context) {
	groupUuid := c.Param("uuid")
	users, err := service.GroupService.GetUserIdByGroupUuid(groupUuid)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, users)
}
