package model

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"im/public"
	"im/public/protocol"
	"time"
)

const MessageLimit = 50 // 最大消息同步数量

// Message 单聊消息
type Message struct {
	ID          int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID      int64     `gorm:"not null;comment:'用户id，指接受者用户id'" json:"user_id"`
	SenderID    int64     `gorm:"not null;comment:'发送者用户id'"`
	SessionType int8      `gorm:"not null;comment:'聊天类型，群聊/单聊'" json:"session_type"`
	ReceiverId  int64     `gorm:"not null;comment:'接收者id，群聊id/用户id'" json:"receiver_id"`
	MessageType int8      `gorm:"not null;comment:'消息类型,语言、文字、图片'" json:"message_type"`
	Content     []byte    `gorm:"not null;comment:'消息内容'" json:"content"`
	Seq         int64     `gorm:"not null;comment:'消息序列号'" json:"seq"`
	SendTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'消息发送时间'" json:"send_time"`
	CreateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (*Message) TableName() string {
	return "message"
}

func ProtoMarshalToMessage(data []byte) []*Message {
	var messages []*Message
	mqMessages := &protocol.MQMessages{}
	err := proto.Unmarshal(data, mqMessages)
	if err != nil {
		zap.S().Error("[Message] [ProtoMarshalToMessage] json.Unmarshal [err] = ", err)
		return nil
	}
	for _, mqMessage := range mqMessages.Messages {
		message := &Message{
			UserID:      mqMessage.UserId,
			SenderID:    mqMessage.SenderId,
			SessionType: int8(mqMessage.SessionType),
			ReceiverId:  mqMessage.ReceiverId,
			MessageType: int8(mqMessage.MessageType),
			Content:     mqMessage.Content,
			Seq:         mqMessage.Seq,
			SendTime:    mqMessage.SendTime.AsTime(),
		}
		messages = append(messages, message)
	}
	return messages
}

func MessageToProtoMarshal(messages ...*Message) []byte {
	if len(messages) == 0 {
		return nil
	}
	var mqMessage []*protocol.MQMessage
	for _, message := range messages {
		mqMessage = append(mqMessage, &protocol.MQMessage{
			UserId:      message.UserID,
			SenderId:    message.SenderID,
			SessionType: int32(message.SessionType),
			ReceiverId:  message.ReceiverId,
			MessageType: int32(message.MessageType),
			Content:     message.Content,
			Seq:         message.Seq,
			SendTime:    timestamppb.New(message.SendTime),
		})
		zap.S().Debug("[Message] [MessageToProtoMarshal] [mqMessage] Time = ", timestamppb.New(message.SendTime))
	}
	bytes, err := proto.Marshal(&protocol.MQMessages{Messages: mqMessage})
	if err != nil {
		zap.S().Error("[Message] [MessageToProtoMarshal] json.Marshal [err] = ", err)
		return nil
	}
	return bytes
}

func MessagesToprotocol(messages []Message) []*protocol.Message {
	protocolMessages := make([]*protocol.Message, 0, len(messages))
	for _, message := range messages {
		protocolMessages = append(protocolMessages, &protocol.Message{
			SessionType: protocol.SessionType(message.SessionType),
			ReceiverId:  message.ReceiverId,
			SenderId:    message.SenderID,
			MessageType: protocol.MessageType(message.MessageType),
			Content:     message.Content,
			Seq:         message.Seq,
		})
	}
	return protocolMessages
}

func CreateMessage(msgs ...*Message) error {
	return public.DB.Create(msgs).Error
}

func ListByUserIdAndSeq(userId, seq int64, limit int) ([]Message, bool, error) {
	var cnt int64
	err := public.DB.Model(&Message{}).Where("user_id = ? and seq > ?", userId, seq).
		Count(&cnt).Error
	if err != nil {
		zap.S().Error("[Message] [ListByUserIdAndSeq] Count [err] = ", err)
		return nil, false, err
	}
	if cnt == 0 {
		zap.S().Debug("[Message] [ListByUserIdAndSeq] Count is 0 | userId = ", userId)
		return nil, false, nil
	}

	var messages []Message
	err = public.DB.Model(&Message{}).Where("user_id = ? and seq > ?", userId, seq).
		Limit(limit).Order("seq ASC").Find(&messages).Error
	if err != nil {
		zap.S().Error("[Message] [ListByUserIdAndSeq] Find [err] = ", err)
		return nil, false, err
	}
	return messages, cnt > int64(limit), nil
}
