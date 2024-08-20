package service

import (
	"im/model"
	"log"
)

// 保存消息
func SaveMessage(message model.Message) {
	// 保存数据库
	result := DbEngine.Create(&message)
	if result.Error != nil {
		log.Println(result)
	}
}

// 获取好友的聊天消息
func GetFriendMessages(userId, friendId, maxId int) []model.Message {
	var messages []model.Message
	DbEngine.Raw("SELECT * FROM messages WHERE id < ? and ((user_id = ? and to_id=?) OR (user_id = ? and to_id=?)) and cmd = 1  ORDER BY id desc LIMIT 10 ", maxId, userId, friendId, friendId, userId).Scan(&messages)
	//
	//log.Printf("SELECT * FROM messages WHERE id < %d and ((user_id = %d and to_id=%d) OR (user_id = %d and to_id=%d))  ORDER BY id desc LIMIT 10 ", maxId, userId, friendId, friendId, userId)

	return messages
}
