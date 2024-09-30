package service

import (
	"go.uber.org/zap"
	"im/internal/model"
	"im/public"

	"github.com/google/uuid"
)

type groupService struct {
}

var GroupService = new(groupService)

func (g *groupService) GetGroups(uuid string) ([]*model.GroupResponse, error) {
	db := public.Db

	queryUser := &model.User{}
	err := db.First(&queryUser, "uuid = ?", uuid).Error
	if err != nil {
		zap.S().Error("GetGroups err = ", err)
		return nil, err
	}

	var groups []*model.GroupResponse
	err = db.Select("g.id AS group_id, g.uuid, g.created_at, g.name, g.notice").
		Model(&model.GroupMember{}).Joins("LEFT JOIN groups AS g ON gm.group_id = g.id").
		Where("gm.user_id = ?", queryUser.Id).Scan(&groups).Error
	if err != nil {
		zap.S().Error("GetGroups err = ", err)
		return nil, err
	}
	return groups, nil
}

func (g *groupService) SaveGroup(userUuid string, group *model.Group) error {
	db := public.Db

	fromUser := &model.User{}
	err := db.Find(&fromUser, "uuid = ?", userUuid).Error
	if err != nil {
		zap.S().Error("SaveGroup err = ", err)
		return err
	}
	group.UserId = fromUser.Id
	group.Uuid = uuid.New().String()
	db.Save(&group)

	groupMember := &model.GroupMember{
		UserId:   fromUser.Id,
		GroupId:  group.ID,
		Nickname: fromUser.Username,
		Mute:     0,
	}
	err = db.Save(groupMember).Error
	if err != nil {
		zap.S().Error("SaveGroup err = ", err)
		return err
	}
	return nil
}

func (g *groupService) GetUserIdByGroupUuid(groupUuid string) ([]*model.User, error) {
	db := public.Db

	group := &model.Group{}
	err := db.First(&group, "uuid = ?", groupUuid).Error
	if err != nil {
		zap.S().Error("GetUserIdByGroupUuid err = ", err)
		return nil, err
	}

	var users []*model.User
	err = db.Select("u.uuid, u.avatar, u.username ").
		Model(&model.Group{}).Joins("JOIN group_members AS gm ON gm.group_id = g.id").
		Joins("JOIN users AS u ON u.id = gm.user_id").
		Where("g.id = ?", group.ID).Scan(&users).Error
	if err != nil {
		zap.S().Error("GetUserIdByGroupUuid err = ", err)
		return nil, err
	}
	return users, nil
}

func (g *groupService) JoinGroup(groupUuid, userUuid string) error {
	db := public.Db

	user := &model.User{}
	err := db.First(&user, "uuid = ?", userUuid).Error
	if err != nil {
		zap.S().Error("JoinGroup err = ", err)
		return err
	}

	group := &model.Group{}
	err = db.First(&group, "uuid = ?", groupUuid).Error
	if err != nil {
		zap.S().Error("JoinGroup err = ", err)
		return err
	}

	groupMember := &model.GroupMember{}
	err = db.First(&groupMember, "user_id = ? and group_id = ?", user.Id, group.ID).Error
	if err != nil {
		zap.S().Error("JoinGroup err = ", err)
		return err
	}

	nickname := user.Nickname
	if nickname == "" {
		nickname = user.Username
	}
	groupMemberInsert := model.GroupMember{
		UserId:   user.Id,
		GroupId:  group.ID,
		Nickname: nickname,
		Mute:     0,
	}
	err = db.Save(&groupMemberInsert).Error
	if err != nil {
		zap.S().Error("JoinGroup err = ", err)
		return err
	}

	return nil
}
