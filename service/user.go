package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/repo"
	"im/util"
)

type LoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type LoginReply struct {
	Token string `json:"token"` // token
}

func Login(c *gin.Context) {
	in := new(LoginRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	// 根据用户名查询用户
	user, err := repo.GetUserByUserName(in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.FailResp(c, "用户名不存在")
			return
		}
		util.FailResp(c, "查询失败")
		return
	}

	if user.Password != util.Md5(in.Password) {
		util.FailResp(c, "密码错误")
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(int64(user.ID))
	if err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, &LoginReply{
		Token: token,
	})
}

type RegisterRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type RegisterReply struct {
	Token string `json:"token"` // token
}

func Register(c *gin.Context) {
	in := new(RegisterRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	// 判断用户名是否存在
	_, err = repo.GetUserByUserName(in.Username)
	if err == nil {
		util.FailResp(c, "用户名已存在")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		util.FailResp(c, "查询失败")
		return
	}

	// 创建账号
	id, err := repo.CreateUser(in.Username, util.Md5(in.Password))
	if err != nil {
		zap.S().Errorf("repo.CreateUser : %v", err)
		util.FailResp(c, "插入失败")
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(int64(id))
	if err != nil {
		zap.S().Errorf("util.GenerateJWT : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, &RegisterReply{
		Token: token,
	})
}
