package public

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"im/config"
)

var Redis *redis.Client

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         config.Configs.Redis.Addr,
		Password:     config.Configs.Redis.Password,
		DB:           config.Configs.Redis.Db,
		PoolSize:     config.Configs.Redis.PoolSize,
		MinIdleConns: config.Configs.Redis.MinIdleConns,
		MaxRetries:   config.Configs.Redis.MaxRetries,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		zap.S().Error("[init] [redisClient.Ping] [err] = ", err.Error())
		panic(err)
	}
	Redis = redisClient
}
