package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"im/config"
	"time"
)

type MyClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func GenerateJWT(Id int64) (string, error) {
	// 创建一个我们自己的声明
	//zap.S().Debugf("时间：%v", time.Now())
	//zap.S().Debugf("有效期时长 %v", time.Duration(config.Configs.Auth.AccessExpire)*time.Second)
	//zap.S().Debugf("token有效期到 %v", time.Now().Add(time.Duration(config.Configs.Auth.AccessExpire)))
	c := MyClaims{
		Id, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.Configs.Auth.AccessExpire) * time.Second).Unix(), // 过期时间
			Issuer:    "zandala",                                                                            // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(config.Configs.Auth.AccessSecret))
}

func VerifyJWT(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Configs.Auth.AccessSecret), nil
	})
	if err != nil {
		zap.S().Errorf("解析token失败: %v", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func GetUid(c *gin.Context) int64 {
	value, exists := c.Get("id")
	if exists {
		zap.S().Info("ctx获取userId失败")
	}
	return value.(int64)
}
