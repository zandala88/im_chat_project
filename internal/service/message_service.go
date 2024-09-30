package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/internal/model"
	"im/internal/websocket/constant"
	"im/internal/websocket/protocol"
	"im/public"
)

type messageService struct {
}

var MessageService = new(messageService)

func (m *messageService) GetMessages(message *model.MessageRequest) ([]*model.MessageResponse, error) {
	db := public.Db

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		queryUser := &model.User{}
		err := db.First(&queryUser, "uuid = ?", message.Uuid).Error

		if err != nil {
			zap.S().Error("GetMessages err = ", err)
			return nil, err
		}

		friend := &model.User{}
		err = db.First(&friend, "username = ?", message.FriendUsername).Error
		if err != nil {
			zap.S().Error("GetMessages err = ", err)
			return nil, err
		}

		var messages []*model.MessageResponse

		err = db.Select("m.id, m.from_user_id, m.to_user_id, m.content, m.content_type, m.url, m.created_at, u.username AS from_username, u.avatar, to_user.username AS to_username").
			Table("messages AS m").
			Joins("LEFT JOIN users AS u ON m.from_user_id = u.id").
			Joins("LEFT JOIN users AS to_user ON m.to_user_id = to_user.id").
			Where("from_user_id IN (?, ?) AND to_user_id IN (?, ?)", queryUser.Id, friend.Id, queryUser.Id, friend.Id).Scan(&messages).Error

		return messages, nil
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		messages, err := fetchGroupMessage(db, message.Uuid)
		if err != nil {
			return nil, err
		}

		return messages, nil
	}

	return nil, errors.New("不支持查询类型")
}

func fetchGroupMessage(db *gorm.DB, toUuid string) ([]*model.MessageResponse, error) {
	group := &model.Group{}
	err := db.First(&group, "uuid = ?", toUuid).Error
	if err != nil {
		zap.S().Error("fetchGroupMessage err = ", err)
		return nil, err
	}

	var messages []*model.MessageResponse

	err = db.Select("m.id, m.from_user_id, m.to_user_id, m.content, m.content_type, m.url, m.created_at, u.username AS from_username, u.avatar").
		Table("messages AS m").
		Joins("LEFT JOIN users AS u ON m.from_user_id = u.id").
		Where("m.message_type = 2 AND m.to_user_id = ?", group.ID).Scan(&messages).Error
	if err != nil {
		zap.S().Error("fetchGroupMessage err = ", err)
		return nil, err
	}

	return messages, nil
}

func (m *messageService) SaveMessage(message protocol.Message) error {
	db := public.Db

	fromUser := &model.User{}
	err := db.Find(&fromUser, "uuid = ?", message.From).Error
	if err != nil {
		zap.S().Error("SaveMessage err = ", err)
		return err
	}

	toUserId := 0
	if message.MessageType == constant.MESSAGE_TYPE_USER {
		toUser := &model.User{}
		err = db.First(&toUser, "uuid = ?", message.To).Error
		if err != nil {
			zap.S().Error("SaveMessage err = ", err)
			return err
		}
		toUserId = toUser.Id
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		group := &model.Group{}
		err = db.Find(&group, "uuid = ?", message.To).Error
		if err != nil {
			zap.S().Error("SaveMessage err = ", err)
			return err
		}
		toUserId = group.ID
	}

	saveMessage := model.Message{
		FromUserId:  fromUser.Id,
		ToUserId:    toUserId,
		Content:     message.Content,
		ContentType: int(message.ContentType),
		MessageType: int(message.MessageType),
		Url:         message.Url,
	}
	err = db.Save(&saveMessage).Error
	if err != nil {
		zap.S().Error("SaveMessage err = ", err)
		return err
	}
	return nil
}
