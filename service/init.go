package service

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"im/model"
	"log"
)

var DbEngine *gorm.DB

func init() {
	dsn := "root:123456@tcp(47.113.118.26:3306)/im?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err.Error())
	}
	DbEngine = db

	db.AutoMigrate(&model.User{}, &model.Community{}, &model.Community{}, &model.Message{}, model.CommunityUsers{})

	fmt.Println(DbEngine)
}
