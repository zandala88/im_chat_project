package config

import (
	"github.com/spf13/viper"
)

var Configs Config

type Config struct {
	MySQL MySQLConfig
	Redis RedisConfig
	Auth  AuthConfig
}

type MySQLConfig struct {
	Port                      string
	Host                      string
	Username                  string
	Password                  string
	Database                  string
	Charset                   string
	ParseTime                 string
	Loc                       string
	IgnoreRecordNotFoundError bool
	LogLevel                  int
	SlowThreshold             int
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}

type RedisConfig struct {
	Addr         string
	Password     string
	Db           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() //根据上面配置加载文件
	if err != nil {
		return
	}
	err = viper.Unmarshal(&Configs)
	if err != nil {
		return
	}
}
