package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// md5 加密 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

// md5 加密 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 对密码加密
func GeneratePasswd(passwd, salt string) string {
	return Md5Encode(passwd + salt)
}

// 校验密码是否正确
func ValidatePasswd(passwd, salt, dbPasswd string) bool {
	return Md5Encode(passwd+salt) == dbPasswd
}
