package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"im/model"
	"im/util"
)

// Register 注册
func Register(c *gin.Context) {
	// 获取参数并验证
	phoneNumber := c.PostForm("phone_number")
	nickname := c.PostForm("nickname")
	password := c.PostForm("password")
	if phoneNumber == "" || password == "" {
		zap.S().Error("[Register] 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询手机号是否已存在
	userRepo := model.NewUserRepo(c)
	cnt, err := userRepo.GetUserCountByPhone(phoneNumber)
	if err != nil {
		zap.S().Error("[Register] [model.GetUserCountByPhone] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}
	if cnt > 0 {
		zap.S().Error("[Register] 账号已被注册")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 插入用户信息
	ub := &model.User{
		PhoneNumber: phoneNumber,
		Nickname:    nickname,
		Password:    util.GetMD5(password),
	}
	err = userRepo.CreateUser(ub)
	if err != nil {
		zap.S().Error("[Register] [model.CreateUser] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(ub.ID)
	if err != nil {
		zap.S().Error("[Register] [util.GenerateJWT(] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 发放 token
	util.SuccessResp(c, gin.H{
		"token":   token,
		"user_id": cast.ToString(ub.ID),
	})
}

// Login 登录
func Login(c *gin.Context) {
	// 验证参数
	phoneNumber := c.PostForm("phone_number")
	password := c.PostForm("password")
	if phoneNumber == "" || password == "" {
		zap.S().Error("[Login] 参数不正确")
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 验证账号名和密码是否正确
	userRepo := model.NewUserRepo(c)
	ub, err := userRepo.GetUserByPhoneAndPassword(phoneNumber, util.GetMD5(password))
	if err != nil {
		zap.S().Error("[Login] [model.GetUserByPhoneAndPassword] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(ub.ID)
	if err != nil {
		zap.S().Error("[Login] [util.GenerateJWT] [err] = ", err.Error())
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, gin.H{
		"token":   token,
		"user_id": cast.ToString(ub.ID),
	})

}
