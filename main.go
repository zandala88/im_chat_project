package main

import (
	_ "im/config"
	"im/public"
	_ "im/public"
	"im/repo"
	"im/router"
)

func main() {
	public.Db.AutoMigrate(&repo.User{})
	router.Router().Run(":8080")
}
