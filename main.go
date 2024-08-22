package main

import (
	_ "im/config"
	_ "im/public"
	"im/router"
)

func main() {
	//public.Db.AutoMigrate(&repo.User{}, &repo.Contact{}, &repo.CommunityUsers{}, &repo.Community{}, &repo.Message{})
	router.Router().Run(":8080")
}
