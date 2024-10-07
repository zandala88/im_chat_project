package model

import (
	"go.uber.org/zap"
	"im/public"
	"time"
)

type User struct {
	ID          int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	PhoneNumber string    `gorm:"not null;unique;comment:'手机号'" json:"phone_number"`
	Nickname    string    `gorm:"not null;comment:'昵称'" json:"nickname"`
	Password    string    `gorm:"not null;comment:'密码'" json:"-"`
	CreateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (*User) TableName() string {
	return "user"
}

func GetUserCountByPhone(phoneNumber string) (int64, error) {
	var cnt int64
	err := public.DB.Model(&User{}).
		Where("phone_number = ?", phoneNumber).Count(&cnt).Error
	if err != nil {
		zap.S().Errorf("GetUserCountByPhone failed, err:%v", err)
		return 0, err
	}
	return cnt, nil
}

func CreateUser(user *User) error {
	err := public.DB.Create(user).Error
	if err != nil {
		zap.S().Errorf("CreateUser failed, err:%v", err)
		return err
	}
	return nil
}

func GetUserByPhoneAndPassword(phoneNumber, password string) (*User, error) {
	user := &User{}
	err := public.DB.Model(&User{}).
		Where("phone_number = ? and password = ?", phoneNumber, password).First(user).Error
	if err != nil {
		zap.S().Errorf("GetUserByPhoneAndPassword failed, err:%v", err)
		return nil, err
	}
	return user, nil
}

func GetUserById(id int64) (*User, error) {
	user := &User{}
	err := public.DB.Model(&User{}).Where("id = ?", id).First(user).Error
	if err != nil {
		zap.S().Errorf("GetUserById failed, err:%v", err)
		return nil, err
	}
	return user, err
}

func GetUserIdByIds(ids []int64) ([]int64, error) {
	var newIds []int64
	m := make(map[int64]struct{}, len(ids))
	for i := 0; i < len(ids); i += 1000 {
		var tmp []int64
		end := i + 1000
		if end > len(ids) {
			end = len(ids)
		}
		subIds := ids[i:end]
		err := public.DB.Model(&User{}).Where("id in ?", subIds).Pluck("id", &tmp).Error
		if err != nil {
			return nil, err
		}
		for _, id := range tmp {
			m[id] = struct{}{}
		}
	}
	for id := range m {
		newIds = append(newIds, id)
	}
	return newIds, nil
}

func GetFriends(userId int64) ([]*User, error) {
	var users []*User
	err := public.DB.Raw("select * from user where id in (select friend_id from friend where user_id = ?)", userId).Scan(&users).Error
	if err != nil {
		zap.S().Errorf("GetFriends failed, err:%v", err)
		return nil, err
	}
	return users, nil
}

func DeleteFriend(userId, friendId int64) error {
	err := public.DB.Where("user_id = ? and friend_id = ?", userId, friendId).
		Or("user_id = ? and friend_id = ?", friendId, userId).
		Delete(&Friend{}).Error
	if err != nil {
		zap.S().Errorf("DeleteFriend failed, err:%v", err)
		return err
	}
	return nil
}
