package model

import (
	"go.uber.org/zap"
	"im/public"
	"time"
)

type Friend struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID     int64     `gorm:"not null;comment:'用户id'" json:"user_id"`
	FriendID   int64     `gorm:"not null;comment:'好友id'" json:"friend_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (*Friend) TableName() string {
	return "friend"
}

func IsFriend(userId, friendId int64) (bool, error) {
	var cnt int64
	err := public.DB.Model(&Friend{}).
		Where("user_id = ? and friend_id = ?", userId, friendId).
		Count(&cnt).Error
	if err != nil {
		zap.S().Errorf("IsFriend failed, err:%v", err)
		return false, err
	}
	return cnt > 0, nil
}

func CreateFriend(friend ...*Friend) error {
	err := public.DB.Create(friend).Error
	if err != nil {
		zap.S().Errorf("CreateFriend failed, err:%v", err)
		return err
	}
	return nil
}
