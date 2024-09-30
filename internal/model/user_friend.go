package model

import (
	"gorm.io/gorm"
	"time"
)

type UserFriend struct {
	ID        int            `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId    int            `json:"userId" gorm:"index;comment:'用户ID'"`
	FriendId  int            `json:"friendId" gorm:"index;comment:'好友ID'"`
}
