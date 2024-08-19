package userLogic

import (
	"github.com/astaxie/beego"
	"go.uber.org/zap"
	"im_api/libs/define"
	"im_api/models/userModel"
	"im_api/module/redis"
	"im_api/module/util"
)

const REDIS_PREFIX = "im_"

type ReturnInfo struct {
	Auth     string
	UserId   string
	UserName string
}

func CheckUserName(userName string) (code int, msg string) {
	status := userModel.CheckoutUserNameExist(userName)
	if status == false {
		code = define.ERR_USER_EXIST_CODE
		msg = define.ERR_USER_EXIST_MSG
	}
	code = define.SUCCESS_CODE
	msg = define.SUCCESS_MSG
	return
}

func CheckAuth(auth string) (code int, msg string, ret ReturnInfo) {
	// var user userModel.User

	err := redis.InitRedis()
	if err != nil {
		zap.S().Errorf("redis init  err: %s", err)
	}

	userInfo, err := redis.HGetAll(auth)
	zap.S().Debugf("userinfo %v", userInfo)
	if err != nil {
		zap.S().Debugf("json err %s", err)
		code = define.ERR_USER_NO_EXIST_CODE
		msg = define.ERR_USER_NO_EXIST_MSG
		return
	}

	if _, ok := userInfo["UserName"]; !ok {
		code = define.ERR_USER_AUTH_INVALID_CODE
		msg = define.ERR_USER_AUTH_INVALID_MSG
		return
	}

	ret.UserName = userInfo["UserName"]
	ret.UserId = userInfo["UserId"]
	ret.Auth = auth

	code = define.SUCCESS_CODE
	msg = define.SUCCESS_MSG
	return

}

func AddOne(user userModel.User) (code int, msg string) {
	user.Id = util.GenUuid()
	user.Password = util.Md5(user.Password)
	code, msg = userModel.AddOne(user)
	return

}

func Login(user userModel.User) (code int, msg string, RetData ReturnInfo) {
	userInfo := userModel.GetUserInfoByUserName(user.UserName)
	// password err
	if util.Md5(user.Password) != userInfo.Password {
		zap.S().Debugf("login password %s, userinfo %s,", user.Password, userInfo.Password)
		code = define.ERR_USER_PASSWORD_CODE
		msg = define.ERR_USER_PASSWORD_MSG
		return
	}
	err := redis.InitRedis()
	if err != nil {
		zap.S().Errorf("redis init  err: %s", err)
	}
	auth := util.GenUuid()
	zap.S().Debugf("user info %v", userInfo)
	userData := make(map[string]interface{}, 2)
	userData["UserId"] = userInfo.Id
	userData["UserName"] = userInfo.UserName

	err = redis.HMSet(auth, userData)

	RetData = ReturnInfo{auth, userInfo.Id, userInfo.UserName}
	if err != nil {
		zap.S().Error("redis set auth err: %s", err)
	}

	code = define.SUCCESS_CODE
	msg = define.SUCCESS_MSG
	return
}

func DeleteAuth(auth string) (code int, msg string) {
	err := redis.InitRedis()
	if err != nil {
		beego.Error("redis init  err: %s", err)
	}

	if err := redis.Delete(auth); err != nil {
		beego.Error("redis delete auth err: %s", err)
	}

	code = define.SUCCESS_CODE
	msg = define.SUCCESS_MSG
	return
}
