package model

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type UserRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewUserRepo(ctx context.Context) *UserRepo {
	return &UserRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (u *UserRepo) GetUserCountByPhone(phoneNumber string) (int64, error) {
	var cnt int64
	err := u.db.Model(&User{}).
		Where("phone_number = ?", phoneNumber).Count(&cnt).Error
	if err != nil {
		zap.S().Error("[User] [GetUserCountByPhone] [err] = ", err)
		return 0, err
	}
	return cnt, nil
}

func (u *UserRepo) CreateUser(user *User) error {
	err := u.db.Create(user).Error
	if err != nil {
		zap.S().Error("[User] [CreateUser] [err] = ", err)
		return err
	}
	return nil
}

func (u *UserRepo) GetUserByPhoneAndPassword(phoneNumber, password string) (*User, error) {
	user := &User{}
	err := u.db.Model(&User{}).
		Where("phone_number = ? and password = ?", phoneNumber, password).First(user).Error
	if err != nil {
		zap.S().Error("[User] [GetUserByPhoneAndPassword] [err] = ", err)
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) GetUserById(id int64) (*User, error) {
	user := &User{}
	err := u.db.Model(&User{}).Where("id = ?", id).First(user).Error
	if err != nil {
		zap.S().Error("[User] [GetUserById] [err] = ", err)
		return nil, err
	}
	return user, err
}

func (u *UserRepo) GetUserIdByIds(ids []int64) ([]int64, error) {
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
			zap.S().Error("[User] [GetUserIdByIds] [err] = ", err)
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

func (u *UserRepo) GetFriends(userId int64) ([]*User, error) {
	var users []*User
	err := u.db.Raw("select * from user where id in (select friend_id from friend where user_id = ?)", userId).Scan(&users).Error
	if err != nil {
		zap.S().Error("[User] [GetFriends] [err] = ", err)
		return nil, err
	}
	return users, nil
}

func (u *UserRepo) DeleteFriend(userId, friendId int64) error {
	err := u.db.Where("user_id = ? and friend_id = ?", userId, friendId).
		Or("user_id = ? and friend_id = ?", friendId, userId).
		Delete(&Friend{}).Error
	if err != nil {
		zap.S().Error("[User] [DeleteFriend] [err] = ", err)
		return err
	}
	return nil
}

func (u *UserRepo) CheckFriendIn(userId int64, friends []int64) (bool, error) {
	var count int64
	err := u.db.Model(&Friend{}).Where("user_id = ?", userId).
		Where("friend_id IN ?", friends).Count(&count).Error
	if err != nil {
		zap.S().Error("[User] [CheckFriendIn] [err] = ", err)
		return false, err
	}
	return len(friends) == int(count), nil
}
