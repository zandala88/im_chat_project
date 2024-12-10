package model

import (
	"context"
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

type GroupRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGroupRepo(ctx context.Context) *GroupRepo {
	return &GroupRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (g *GroupRepo) CreateGroup(group *Group, ids []int64) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
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

func (g *GroupRepo) GetGroupById(groupId int64) (*Group, error) {
	var group = &Group{}
	err := g.db.First(group, groupId).Error
	if err != nil {
		zap.S().Error("[Group] [GetGroupById] [err] = ", err)
		return nil, err
	}
	return group, nil
}

func (g *GroupRepo) GetGroups(userId int64) ([]*Group, error) {
	var groups = make([]*Group, 0)
	err := g.db.Table("group").
		Joins("join group_user on group.id = group_user.group_id").
		Where("group_user.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		zap.S().Error("[Group] [GetGroups] [err] = ", err)
		return nil, err
	}
	return groups, nil
}

func (g *GroupRepo) IsGroupOwner(userId, groupId int64) (bool, error) {
	var cnt int64
	err := g.db.Model(&Group{}).
		Where("owner_id = ? and id = ?", userId, groupId).
		Count(&cnt).Error
	if err != nil {
		zap.S().Error("[Group] [IsGroupOwner] [err] = ", err)
		return false, err
	}
	return cnt > 0, nil
}

func (g *GroupRepo) DeleteGroup(groupId int64) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
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
