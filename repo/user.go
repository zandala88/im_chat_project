package repo

import (
	"im/public"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	UserName  string    `gorm:"column:user_name;type:varchar(50);" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(36);" json:"password"`
	Mobile    string    `gorm:"column:mobile;type:varchar(20);" json:"mobile"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}

func GetUserByUserName(username string) (*User, error) {
	db := public.Db
	data := &User{}
	err := db.Where("username", username).First(&data).Error
	return data, err
}

func GetUserByUserId(userId int64) (*User, error) {
	db := public.Db
	data := &User{}
	err := db.Where("id", userId).First(&data).Error
	return data, err
}

func CreateUser(username, password string) (int64, error) {
	db := public.Db
	data := &User{
		UserName: username,
		Password: password,
	}
	err := db.Create(&data).Error
	return data.Id, err
}
