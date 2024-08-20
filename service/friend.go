package service

import (
	"errors"
	"im/model"
)

// 添加好友
func AddFriend(userId, friendId int) (model.Contact, error) {
	// 首先查找好友是否存在
	contact := model.Contact{}
	if userId == friendId {
		return contact, errors.New("自己不能添加自己为好友")
	}

	// 查找好友是否存在
	user := model.User{ID: userId}
	DbEngine.Find(&user)
	if user.ID == 0 {
		return contact, errors.New("好友不存在")
	}
	DbEngine.Where("user_id = ?", userId).Where("friend_id = ?", friendId).Find(&contact)
	if contact.ID > 0 {
		return contact, errors.New("好友已经存在")
	}

	// 如果找到了, 则更新记录
	contact.UserId = userId
	contact.FriendId = friendId
	DbEngine.Begin()
	result := DbEngine.Create(&contact)
	result2 := DbEngine.Create(&model.Contact{
		UserId:   friendId,
		FriendId: userId,
	})
	if result.Error != nil || result2.Error != nil {
		DbEngine.Rollback()
		if result.Error != nil {
			return contact, result.Error
		}

		return contact, result2.Error
	}
	DbEngine.Commit()

	return contact, nil
}

// 删除好友
func DeleteFriend(userId, friendId int) (ok bool, error error) {
	// 查找是否好友
	contact := model.Contact{}
	DbEngine.Where("user_id = ?", userId).Where("friend_id = ?", friendId).First(&contact)
	if contact.ID == 0 {
		return false, errors.New("此用户并不是好友")
	}
	DbEngine.Begin()
	// 删除彼此的好友关系
	result1 := DbEngine.Delete(&model.Contact{
		UserId:   userId,
		FriendId: friendId,
	})
	result2 := DbEngine.Delete(&model.Contact{
		UserId:   friendId,
		FriendId: userId,
	})

	if result1.Error != nil || result2.Error != nil {
		DbEngine.Rollback()
		if result1.Error != nil {
			return false, result1.Error
		}
		return false, result2.Error
	}
	DbEngine.Commit()

	return true, nil
}

// 获取好友列表
func Friends(user model.User) ([]model.User, error) {
	contacts := []model.Contact{}
	// 好友 ids
	friendIds := []int{}
	DbEngine.Where("user_id = ?", user.ID).Select("friend_id").Find(&contacts)
	for _, val := range contacts {
		friendIds = append(friendIds, int(val.FriendId))
	}

	users := []model.User{}
	if len(friendIds) == 0 {
		return users, nil
	}

	// 查找要有用户
	DbEngine.Find(&users, friendIds)

	return users, nil
}

// 获取好友信息
func GetFriend(friendId int) (user model.User, err error) {
	user = model.User{}
	DbEngine.Where("id = ? ", friendId).Find(&user)
	if user.ID == 0 {
		return user, errors.New("未找到好友信息")
	}

	return user, nil
}
