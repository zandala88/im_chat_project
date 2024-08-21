package main

import (
	_ "im/config"
	_ "im/public"
	"im/router"
)

func main() {
	//public.Db.AutoMigrate(&repo.User{}, &repo.Contact{})
	router.Router().Run(":8080")
}
