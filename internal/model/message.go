package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID          int            `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
	FromUserId  int            `json:"fromUserId" gorm:"index"`
	ToUserId    int            `json:"toUserId" gorm:"index;comment:'发送给端的id，可为用户id或者群id'"`
	Content     string         `json:"content" gorm:"type:varchar(2500)"`
	MessageType int            `json:"messageType" gorm:"comment:'消息类型：1单聊，2群聊'"`
	ContentType int            `json:"contentType" gorm:"comment:'消息内容类型：1文字 2.普通文件 3.图片 4.音频 5.视频 6.语音聊天 7.视频聊天'"`
	Pic         string         `json:"pic" gorm:"type:text;comment:'缩略图"`
	Url         string         `json:"url" gorm:"type:varchar(350);comment:'文件或者图片地址'"`
}
