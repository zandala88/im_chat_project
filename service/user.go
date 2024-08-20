package service

import (
	"errors"
	"fmt"
	"im/model"
	"im/util"
	"log"
	"math/rand"
	"time"
)

type UserService struct{}

// 注册函数
func (us UserService) Register(mobile, passwd, nickname string) (model.User, error) {
	// 检测手机号码是否存在, 存在提示已经注册, 不存在数据可插入操作
	user := model.User{}
	DbEngine.Where("mobile = ?", mobile).First(&user)
	if user.ID > 0 {
		return user, errors.New("该手机号已经注册")
	}

	// 插入用户数据
	user.Mobile = mobile
	user.NickName = nickname
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	user.Password = util.GeneratePasswd(passwd, user.Salt)

	result := DbEngine.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// 用户登录
func (us UserService) Login(mobile, passwd string) (model.User, error) {
	user := model.User{}
	DbEngine.Where("mobile = ?", mobile).First(&user)
	if user.ID == 0 {
		return user, errors.New("用户不存在")
	}

	isValidate := util.ValidatePasswd(passwd, user.Salt, user.Password)
	if !isValidate {
		return user, errors.New("用户密码不正确")
	}
	token := util.MD5Encode(fmt.Sprintf("%d", time.Now().UnixNano()))
	// 更新用户 token
	result := DbEngine.Model(&user).Update("token", token)
	if result.Error != nil {
		log.Println(result.Error)
	}
	user.Token = token
	// 验证密码
	return user, nil
}
