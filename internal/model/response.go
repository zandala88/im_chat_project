package model

import "time"

type MessageResponse struct {
	ID           int       `json:"id" gorm:"primarykey"`
	FromUserId   int       `json:"fromUserId" gorm:"index"`
	ToUserId     int       `json:"toUserId" gorm:"index"`
	Content      string    `json:"content" gorm:"type:varchar(2500)"`
	ContentType  int       `json:"contentType" gorm:"comment:'消息内容类型：1文字，2语音，3视频'"`
	CreatedAt    time.Time `json:"createAt"`
	FromUsername string    `json:"fromUsername"`
	ToUsername   string    `json:"toUsername"`
	Avatar       string    `json:"avatar"`
	Url          string    `json:"url"`
}

type GroupResponse struct {
	Uuid      string    `json:"uuid"`
	GroupId   int       `json:"groupId"`
	CreatedAt time.Time `json:"createAt"`
	Name      string    `json:"name"`
	Notice    string    `json:"notice"`
}

type SearchResponse struct {
	User  *User  `json:"user"`
	Group *Group `json:"group"`
}
