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
			zap.S().Errorf("CreateGroup failed, err:%v", err)
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
