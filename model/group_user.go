package model

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type GroupUserRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGroupUserRepo(ctx context.Context) *GroupUserRepo {
	return &GroupUserRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

// IsBelongToGroup 验证用户是否属于群
func (g *GroupUserRepo) IsBelongToGroup(userId, groupId int64) (bool, error) {
	var cnt int64
	err := g.db.Model(&GroupUser{}).
		Where("user_id = ? and group_id = ?", userId, groupId).
		Count(&cnt).Error
	if err != nil {
		zap.S().Error("[GroupUser] [IsBelongToGroup] [err] = ", err)
		return false, err
	}
	return cnt > 0, nil
}

func (g *GroupUserRepo) GetGroupUserIdsByGroupId(groupId int64) ([]int64, error) {
	var ids []int64
	err := g.db.Model(&GroupUser{}).
		Where("group_id = ?", groupId).Pluck("user_id", &ids).Error
	if err != nil {
		zap.S().Error("[GroupUser] [GetGroupUserIdsByGroupId] [err] = ", err)
		return nil, err
	}
	return ids, nil
}

func (g *GroupUserRepo) JoinGroup(groupId, userId int64) error {
	err := g.db.Create(&GroupUser{
		GroupID: groupId,
		UserID:  userId,
	}).Error
	if err != nil {
		zap.S().Error("[GroupUser] [JoinGroup] [err] = ", err)
		return err
	}
	return nil
}

func (g *GroupUserRepo) ExitGroup(groupId, userId int64) error {
	err := g.db.Where("group_id = ? and user_id = ?", groupId, userId).Delete(&GroupUser{}).Error
	if err != nil {
		zap.S().Error("[GroupUser] [ExitGroup] [err] = ", err)
		return err
	}
	return nil
}
