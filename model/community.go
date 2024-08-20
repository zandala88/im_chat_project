package model

import "time"

//CREATE TABLE user_group (
//id INT(10) NOT NULL PRIMARY KEY auto_increment,
//name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '群名称',
//avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '群图片',
//owner_id INT(10) NOT NULL DEFAULT 0 COMMENT '群创建人 id',
//created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
//updated_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
//deleted_at TIMESTAMP NULL
//) COMMENT '群表';

type Community struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"varchar(50); not null; default:''" json:"name"`
	OwnerId   int       `gorm:"int(10); not null; default:0" json:"owner_id"`
	Avatar    string    `gorm:"varchar(50); not null; default:''" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"default:null" json:"deleted_at"`
}

//CREATE TABLE community_users (
//id int(10) NOT NULL PRIMARY KEY auto_increment,
//community_id int(10) NOT NULL DEFAULT 0,
//user_id INT(10) NOT NULL DEFAULT 0
//) COMMENT '群组用户';

// 用户与群组的关系
type CommunityUsers struct {
	ID          int `gorm: "primary"; json:"id"`
	CommunityId int `json:"community_id"`
	UserId      int `json:"user_id"`
}
