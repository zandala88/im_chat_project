package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

func GenUuid() string {
	uuidStr := uuid.Must(uuid.NewV4(), nil).String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

func JsonEncode(structModel interface{}) (string, error) {
	jsonStr, err := json.Marshal(structModel)
	return string(jsonStr), err
}
func JsonDecode(jsonStr string, structModel interface{}) error {
	decode := json.NewDecoder(strings.NewReader(jsonStr))
	err := decode.Decode(structModel)
	return err
}

func Md5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}
