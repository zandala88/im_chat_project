package model

type MessageRequest struct {
	MessageType    int    `json:"messageType" form:"messageType"`
	Uuid           string `json:"uuid" form:"uuid"`
	FriendUsername string `json:"friendUsername" form:"friendUsername"`
}

type FriendRequest struct {
	Uuid           string `json:"uuid"`
	FriendUsername string `json:"friendUsername"`
}

type ModifyUserInfoRequest struct {
	Uuid     string `json:"uuid" binding:"required" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`
	Nickname string `json:"nickname" gorm:"comment:'昵称'"`
	Avatar   string `json:"avatar" gorm:"type:varchar(150);comment:'头像'"`
	Email    string `json:"email" gorm:"type:varchar(80);column:email;comment:'邮箱'"`
	Password string `json:"password" form:"password" gorm:"type:varchar(150);not null; comment:'密码'"`
}
