package model

//CREATE TABLE contact (
//id INT(10) NOT NULL PRIMARY key auto_increment,
//user_id int(10) not null DEFAULT 0 COMMENT '用户 ID',
//friend_id int(10) not null default 0 COMMENT '好友 ID',
//created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//deleted_at TIMESTAMP NULL,
//index user_id_index (user_id),
//index friend_id_index (friend_id)
//) COMMENT '用户好友表';

type Contact struct {
	ID       int
	UserId   int `gorm:"int(10)" json:"user_id"`
	FriendId int `gorm:"int(10)" json:"friend_id"`
}
