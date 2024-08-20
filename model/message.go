package model

import "time"

const (
	MEDIA_TYPE_TEXT = iota + 1
	// 新闻格式
	MEDIA_TYPE_NEWS
	// 语音格式
	MEDIA_TYPE_VOICE
	// 图片格式
	MEDIA_TYPE_IMG
	// 红包格式
	MEDIA_TYPE_REDPACKAGE
	// emoj 表情信息
	MEDIA_TYPE_EMOJ
	// 链接格式
	MEDIA_TYPE_LINK
	// 视频格式
	MEDIA_TYPE_VIDEO
	// 名片格式
	MEDIA_TYPE_CONCAT
	// 其他格式, 前端做相应的处理
	MEDIA_TYPE_UNDE = 100
)

const (
	CMD_PRIVATE = iota + 1 // 私聊
	CMD_GROUP              // 群聊
)

type Message struct {
	ID        int       `json:"id,omitempty" form:"id"`           // 消息 id
	UserId    uint      `json:"user_id,omitempty" form:"userid"`  // 发送者用户 id
	Cmd       int       `json:"cmd,omitempty" form:"cmd"`         // 群聊还是私聊
	ToId      int       `json:"to_id omitempty" form:"to_id"`     //对端 id 或者群聊 id
	Media     uint      `json:"media,omitempty" form:"media"`     // 消息样式
	Content   string    `json:"content,omitempty" form:"content"` // 消息的内容
	Pic       string    `json:"pic,omitempty" form:"pic"`         // 图片预览
	Url       string    `json:"url,omitempty" form:"url"`         // 服务的 url
	Memo      string    `json:"memo,omitempty" form:"memo"`       // 简单描述
	Amount    int       `json:"amount,omitempty" form:"amount"`   // 和数字相关
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}
