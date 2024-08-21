package repo

import (
	"gorm.io/gorm"
	"im/public"
)

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password string `gorm:"column:password;type:varchar(36);" json:"password"`
	Mobile   string `gorm:"column:mobile;type:varchar(20);" json:"mobile"`
}

func GetUserByUserName(username string) (*User, error) {
	db := public.Db
	data := &User{}
	err := db.Where("username", username).First(&data).Error
	return data, err
}

func CreateUser(username, password string) (uint, error) {
	db := public.Db
	data := &User{
		Username: username,
		Password: password,
	}
	err := db.Create(&data).Error
	return data.ID, err
}
