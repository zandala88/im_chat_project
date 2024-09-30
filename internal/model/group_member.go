package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupMember struct {
	ID        int            `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId    int            `json:"userId" gorm:"index;comment:'用户ID'"`
	GroupId   int            `json:"groupId" gorm:"index;comment:'群组ID'"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(350);comment:'昵称"`
	Mute      int            `json:"mute" gorm:"comment:'是否禁言'"`
}
