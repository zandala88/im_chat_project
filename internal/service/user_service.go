package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/internal/model"
	"im/public"
	"time"

	"github.com/google/uuid"
)

type userService struct{}

var UserService = new(userService)

func (u *userService) Register(user *model.User) error {
	db := public.Db

	var userCount int64
	err := db.Model(user).Where("username", user.Username).Count(&userCount).Error
	if err != nil {
		zap.S().Error("Register err = ", err)
		return err
	}
	if userCount > 0 {
		return errors.New("user already exists")
	}
	user.Uuid = uuid.New().String()
	user.CreateAt = time.Now()
	user.DeleteAt = 0

	if err = db.Create(&user).Error; err != nil {
		zap.S().Error("Register err = ", err)
		return err
	}
	return nil
}

func (u *userService) Login(user *model.User) (bool, error) {
	db := public.Db
	zap.S().Debugf("Login %#v", user)

	queryUser := &model.User{}
	err := db.First(&queryUser, "username = ?", user.Username).Error
	if err != nil {
		zap.S().Error("Login err = ", err)
		return false, err
	}
	zap.S().Debugf("Login %#v", queryUser)

	user.Uuid = queryUser.Uuid

	return queryUser.Password == user.Password, nil
}

func (u *userService) ModifyUserInfo(user *model.ModifyUserInfoRequest) error {
	db := public.Db

	queryUser := &model.User{}
	err := db.First(&queryUser, "uuid = ?", user.Uuid).Error
	if err != nil {
		zap.S().Error("ModifyUserInfo err = ", err)
		return err
	}

	queryUser.Nickname = user.Nickname
	queryUser.Email = user.Email
	queryUser.Password = user.Password

	db.Updates(queryUser)
	return nil
}

func (u *userService) GetUserDetails(uuid string) (*model.User, error) {
	db := public.Db

	queryUser := &model.User{}
	err := db.Select("uuid", "username", "nickname", "avatar").
		First(&queryUser, "uuid = ?", uuid).Error
	if err != nil {
		zap.S().Error("GetUserDetails err = ", err)
		return nil, err
	}

	return queryUser, nil
}

func (u *userService) GetUserOrGroupByName(name string) (*model.SearchResponse, error) {
	db := public.Db

	queryUser := &model.User{}
	err := db.Select("uuid", "username", "nickname", "avatar").First(&queryUser, "username = ?", name).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Error("GetUserOrGroupByName err = ", err)
		return nil, err
	}

	queryGroup := &model.Group{}
	err = db.Select("uuid", "name").First(&queryGroup, "name = ?", name).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Error("GetUserOrGroupByName err = ", err)
		return nil, err
	}

	return &model.SearchResponse{
		User:  queryUser,
		Group: queryGroup,
	}, nil
}

func (u *userService) GetUserList(uuid string) ([]*model.User, error) {
	db := public.Db

	queryUser := &model.User{}
	err := db.First(&queryUser, "uuid = ?", uuid).Error
	if err != nil {
		zap.S().Error("GetUserList err = ", err)
		return nil, err
	}

	var queryUsers []*model.User
	err = db.Select("u.username, u.uuid, u.avatar").
		Table("user_friends AS uf").Joins("JOIN users AS u ON uf.friend_id = u.id").
		Where("uf.user_id = ?", queryUser.Id).
		Scan(&queryUsers).Error

	if err != nil {
		zap.S().Error("GetUserList err = ", err)
		return nil, err
	}

	return queryUsers, nil
}

func (u *userService) AddFriend(userFriendRequest *model.FriendRequest) error {
	db := public.Db

	queryUser := &model.User{}
	err := db.First(&queryUser, "uuid = ?", userFriendRequest.Uuid).Error
	if err != nil {
		zap.S().Error("AddFriend err = ", err)
		return err
	}
	zap.S().Debugf("AddFriend %#v", queryUser)

	friend := &model.User{}
	err = db.First(&friend, "username = ?", userFriendRequest.FriendUsername).Error
	if err != nil {
		zap.S().Error("AddFriend err = ", err)
		return err
	}

	userFriend := &model.UserFriend{
		UserId:   queryUser.Id,
		FriendId: friend.Id,
	}

	userFriendQuery := &model.UserFriend{}
	err = db.First(&userFriendQuery, "user_id = ? and friend_id = ?", queryUser.Id, friend.Id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Error("AddFriend err = ", err)
		return err
	}

	err = db.Save(&userFriend).Error
	if err != nil {
		zap.S().Error("AddFriend err = ", err)
		return err
	}

	zap.S().Debugf("userFriend %#v", userFriend)
	return nil
}

func (u *userService) ModifyUserAvatar(avatar string, userUuid string) error {
	db := public.Db

	queryUser := &model.User{}
	err := db.First(&queryUser, "uuid = ?", userUuid).Error
	if err != nil {
		zap.S().Error("ModifyUserAvatar err = ", err)
	}

	err = db.Model(&queryUser).Update("avatar", avatar).Error
	if err != nil {
		zap.S().Error("ModifyUserAvatar err = ", err)
	}
	return nil
}
