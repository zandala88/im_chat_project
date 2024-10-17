package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/public"
	"time"
)

type Group struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	Name       string    `gorm:"not null;comment:'群组名称'" json:"name"`
	OwnerID    int64     `gorm:"not null;comment:'群主id'" json:"owner_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (*Group) TableName() string {
	return "group"
}

func CreateGroup(group *Group, ids []int64) error {
	return public.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(group).Error
		if err != nil {
			zap.S().Error("[Group] [CreateGroup] [err] = ", err)
			return err
		}

		groupUsers := make([]*GroupUser, 0, len(ids))
		for _, id := range ids {
			groupUsers = append(groupUsers, &GroupUser{
				GroupID: group.ID,
				UserID:  id,
			})
		}
		return tx.Create(groupUsers).Error
	})
}

func GetGroupById(groupId int64) (*Group, error) {
	group := new(Group)
	err := public.DB.First(group, groupId).Error
	if err != nil {
		zap.S().Error("[Group] [GetGroupById] [err] = ", err)
		return nil, err
	}
	return group, nil
}

func GetGroups(userId int64) ([]*Group, error) {
	groups := make([]*Group, 0)
	err := public.DB.Table("group").
		Joins("join group_user on group.id = group_user.group_id").
		Where("group_user.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		zap.S().Error("[Group] [GetGroups] [err] = ", err)
		return nil, err
	}
	return groups, nil
}

func IsGroupOwner(userId, groupId int64) (bool, error) {
	var cnt int64
	err := public.DB.Model(&Group{}).
		Where("owner_id = ? and id = ?", userId, groupId).
		Count(&cnt).Error
	if err != nil {
		zap.S().Error("[Group] [IsGroupOwner] [err] = ", err)
		return false, err
	}
	return cnt > 0, nil
}

func DeleteGroup(groupId int64) error {
	return public.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("group_id = ?", groupId).Delete(&GroupUser{}).Error
		if err != nil {
			zap.S().Error("[Group] [DeleteGroup] [err] = ", err)
			return err
		}

		err = tx.Delete(&Group{}, groupId).Error
		if err != nil {
			zap.S().Error("[Group] [DeleteGroup] [err] = ", err)
			return err
		}
		return nil
	})
}
