package util

import (
	"crypto/md5"
	"fmt"
	"im/config"
)

func GetMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s+config.Configs.App.Salt)))
}
