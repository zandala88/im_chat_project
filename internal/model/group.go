package model

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID        int            `json:"id" gorm:"primarykey"`
	Uuid      string         `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`
	CreatedAt time.Time      `json:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId    int            `json:"userId" gorm:"index;comment:'群主ID'"`
	Name      string         `json:"name" gorm:"type:varchar(150);comment:'群名称"`
	Notice    string         `json:"notice" gorm:"type:varchar(350);comment:'群公告"`
}
