package public

import (
	"context"
	"github.com/go-redis/redis/v8"
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
		panic(err)
	}
	Redis = redisClient
}
