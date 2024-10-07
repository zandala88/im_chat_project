package model

import (
	"go.uber.org/zap"
	"im/public"
	"time"
)

type GroupUser struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	GroupID    int64     `gorm:"not null;comment:'组id'" json:"group_id"`
	UserID     int64     `gorm:"not null;comment:'用户id'" json:"user_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (*GroupUser) TableName() string {
	return "group_user"
}

// IsBelongToGroup 验证用户是否属于群
func IsBelongToGroup(userId, groupId int64) (bool, error) {
	var cnt int64
	err := public.DB.Model(&GroupUser{}).
		Where("user_id = ? and group_id = ?", userId, groupId).
		Count(&cnt).Error
	if err != nil {
		zap.S().Errorf("IsBelongToGroup failed, err:%v", err)
		return false, err
	}
	return cnt > 0, nil
}

func GetGroupUserIdsByGroupId(groupId int64) ([]int64, error) {
	var ids []int64
	err := public.DB.Model(&GroupUser{}).
		Where("group_id = ?", groupId).Pluck("user_id", &ids).Error
	if err != nil {
		zap.S().Errorf("GetGroupUserIdsByGroupId failed, err:%v", err)
		return nil, err
	}
	return ids, nil
}
