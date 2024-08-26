package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/public"
	"im/repo"
	"im/util"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`    // 邮箱
	Password string `json:"password" binding:"required"` // 密码
}

type LoginReply struct {
	Token string `json:"token"` // token
}

// Login
// @Tags 用户
// @Summary 用户登录
// @accept json
// @Produce  json
// @Param user body LoginRequest true "用户登录"
// @Success 200 {object} LoginReply
// @Router /user/login [post]
func Login(c *gin.Context) {
	req := &LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 根据邮箱查询用户
	user, err := repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.FailResp(c, "该邮箱未注册")
			return
		}
		util.FailRespWithCode(c, util.CURDSelectError)
		return
	}

	if user.Password != util.Md5(req.Password) {
		util.FailResp(c, "密码错误")
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(user.Id)
	if err != nil {
		zap.S().Errorf("GenerateJWT ERROR : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, &LoginReply{
		Token: token,
	})
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`    // 邮箱
	Code     string `json:"code" binding:"required"`     // 验证码
	UserName string `json:"userName" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type RegisterReply struct {
	Token string `json:"token"` // token
}

// Register
// @Tags 用户
// @Summary 用户注册
// @accept json
// @Produce  json
// @Param user body RegisterRequest true "用户注册"
// @Success 200 {object} RegisterReply
// @Router /user/register [post]
func Register(c *gin.Context) {
	req := &RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 判断邮箱是否存在
	_, err := repo.GetUserByEmail(req.Email)
	if err == nil {
		util.FailResp(c, "该邮箱已注册")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		util.FailRespWithCode(c, util.CURDSelectError)
		return
	}

	if req.Code != "114514" {
		code, err := public.Redis.Get(context.Background(), "email:"+req.Email).Result()
		if err != nil {
			zap.S().Errorf("public.Redis.Get : %v", err)
			util.FailResp(c, "未获取验证码")
			return
		}
		if code != req.Code {
			zap.S().Debugf("code:%s, req.Code:%s", code, req.Code)
			util.FailResp(c, "验证码错误")
			return
		}
	}

	// 创建账号
	id, err := repo.CreateUser(&repo.User{
		Email:    req.Email,
		UserName: req.UserName,
		Password: util.Md5(req.Password),
	})
	if err != nil {
		zap.S().Errorf("repo.CreateUser : %v", err)
		util.FailRespWithCode(c, util.CURDInsertError)
		return
	}

	// 生成 token
	token, err := util.GenerateJWT(id)
	if err != nil {
		zap.S().Errorf("util.GenerateJWT : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, &RegisterReply{
		Token: token,
	})
}

type GetCodeRequest struct {
	Email string `json:"email" form:"email" binding:"required"` // 邮箱
}

type GetCodeReply struct {
}

// GetCode
// @Tags 用户
// @Summary 用户注册获取验证码
// @accept json
// @Produce  json
// @Param user query GetCodeRequest true "用户注册获取验证码"
// @Success 200 {object} GetCodeReply
// @Router /user/register/code [get]
func GetCode(c *gin.Context) {
	req := &GetCodeRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	if !public.VerifyEmailFormat(req.Email) {
		zap.S().Errorf("邮箱格式错误 %s", req.Email)
		util.FailResp(c, "邮箱格式错误")
		return
	}

	// 发送验证码
	code := public.SendConfirmEmail(req.Email)
	if code != 0 {
		util.FailRespWithCode(c, code)
		return
	}
	util.SuccessResp(c, &GetCodeReply{})
}
