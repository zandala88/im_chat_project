package repo

import (
	"gorm.io/gorm"
	"im/public"
	"time"
)

type Contact struct {
	Id        int64     `json:"id"`
	UserId    int64     `gorm:"column:user_id;type:int(10)" json:"user_id"`
	FriendId  int64     `gorm:"column:friend_id;type:int(10)" json:"friend_id"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}

func GetOneContact(userId, friendId int64) (*Contact, error) {
	db := public.Db
	data := &Contact{}
	err := db.Where("user_id = ? AND friend_id = ?", userId, friendId).First(&data).Error
	return data, err
}

func AddFriendCreate(userId, friendId int64) error {
	db := public.Db
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(&Contact{UserId: userId, FriendId: friendId}).Error; err != nil {
			return err
		}
		if err := db.Create(&Contact{UserId: friendId, FriendId: userId}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetFriendList(userId int64) ([]int64, error) {
	db := public.Db
	var data []int64
	err := db.Select("friend_id").Where("user_id", userId).Scan(&data).Error
	return data, err
}
