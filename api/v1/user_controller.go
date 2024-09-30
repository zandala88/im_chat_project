package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"im/internal/model"
	"im/internal/service"
	"im/util"
)

func Login(c *gin.Context) {
	user := &model.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		zap.S().Error("ShouldBindJSON err = ", err)
		util.FailResp(c, err.Error())
		return
	}
	zap.S().Debugf("user %#v", user)

	ok, err := service.UserService.Login(user)
	if err != nil || !ok {
		util.FailResp(c, "Login failed")
		return
	}
	util.SuccessResp(c, user)

}

func Register(c *gin.Context) {
	user := &model.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		zap.S().Error("ShouldBindJSON err = ", err)
		util.FailResp(c, err.Error())
		return
	}

	err = service.UserService.Register(user)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, user)
}

func ModifyUserInfo(c *gin.Context) {
	user := &model.ModifyUserInfoRequest{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		zap.S().Error("ShouldBindJSON err = ", err)
		util.FailResp(c, err.Error())
		return
	}
	zap.S().Debugf("user %#v", user)
	if err := service.UserService.ModifyUserInfo(user); err != nil {
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, nil)
}

func GetUserDetails(c *gin.Context) {
	uuid := c.Param("uuid")
	details, err := service.UserService.GetUserDetails(uuid)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, details)
}

func GetUserOrGroupByName(c *gin.Context) {
	name := c.Query("name")
	byName, err := service.UserService.GetUserOrGroupByName(name)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, byName)
}

func GetUserList(c *gin.Context) {
	uuid := c.Query("uuid")
	list, err := service.UserService.GetUserList(uuid)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	util.SuccessResp(c, list)
}

func AddFriend(c *gin.Context) {
	userFriendRequest := &model.FriendRequest{}
	err := c.ShouldBindJSON(&userFriendRequest)
	if err != nil {
		zap.S().Error("ShouldBindJSON err = ", err)
		util.FailResp(c, err.Error())
		return
	}
	zap.S().Info("userFriendRequest %#v", userFriendRequest)

	err = service.UserService.AddFriend(userFriendRequest)
	if nil != err {
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, nil)
}
