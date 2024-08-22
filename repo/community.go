package repo

import (
	"im/public"
	"time"
)

type Community struct {
	Id        int64     `json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50); not null; default:''" json:"name"`
	OwnerId   int64     `gorm:"column:owner_id;type:int(10); not null; default:0" json:"owner_id"`
	Avatar    string    `gorm:"column:avatar;type:varchar(50); not null; default:''" json:"avatar"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}

// CommunityUsers 用户与群组的关系
type CommunityUsers struct {
	Id          int64 `json:"id"`
	CommunityId int64 `gorm:"column:community_id" json:"community_id"`
	UserId      int64 `gorm:"column:user_id" json:"user_id"`
}

// GetCommunitiesByUserId 获取用户的群组信息
func GetCommunitiesByUserId(userId int) ([]*CommunityUsers, error) {
	db := public.Db
	var communityUsers []*CommunityUsers
	err := db.Where("user_id = ?", userId).Find(&communityUsers).Error
	return communityUsers, err
}
