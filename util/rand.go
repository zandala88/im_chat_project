package util

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.NewV4().String()
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
