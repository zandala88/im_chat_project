package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Mobile    string    `gorm:"type:varchar(11);unique;default:'';" json:"mobile"`
	Password  string    `gorm:"type:varchar(40)" json:"-"`
	Avatar    string    `gorm:"type:varchar(160)" json:"avatar"`
	Sex       int       `json:"sex"`
	NickName  string    `gorm:"type:varchar(20)" json:"nick_name"`
	Salt      string    `gorm:"type:varchar(10)" json:"-"`
	Online    int8      `json:"online"`
	Token     string    `gorm:"type:varchar(50)" json:"token"`
	Memo      string    `gorm:"type:varchar(255)" json:"memo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}
