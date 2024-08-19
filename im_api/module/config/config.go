package config

import (
	"github.com/astaxie/beego/config"
	"os"
)

func Reader(filePath string) (iniconf config.Configer, err error) {
	path, _ := os.Getwd()
	// path += "/../.."
	appconf, err := config.NewConfig("ini", path+"/conf/app.conf")
	if err != nil {
		return
	}
	runmode := appconf.String("runmode")
	if runmode == "dev" {
		iniconf, err = config.NewConfig("ini", path+"/conf/dev/"+filePath)
	} else {
		iniconf, err = config.NewConfig("ini", path+"/conf/prod/"+filePath)
	}
	return iniconf, err

}
