package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"im/config"
	"time"
)

type MyClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func GenerateJWT(Id int64) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		Id, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.Configs.Auth.AccessExpire)).Unix(), // 过期时间
			Issuer:    "zandala",                                                              // 签发人
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
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
