package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"im/public"
	"time"
)

const (
	userOnlinePrefix = "user_online_" // 用户在线状态设置
	ttl1D            = 24 * 60 * 60   // s  1天
)

func getUserKey(userId int64) string {
	return fmt.Sprintf("%s%d", userOnlinePrefix, userId)
}

// SetUserOnline 设置用户在线
func SetUserOnline(userId int64, addr string) error {
	key := getUserKey(userId)
	_, err := public.Redis.Set(context.Background(), key, addr, ttl1D*time.Second).Result()
	if err != nil {
		zap.S().Error("[设置用户在线] 错误, err:", err)
		return err
	}
	return nil
}

// GetUserOnline 获取用户在线地址
// 如果获取不到，返回 addr = "" 且 err 为 nil
func GetUserOnline(userId int64) (string, error) {
	key := getUserKey(userId)
	addr, err := public.Redis.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.S().Error("[获取用户在线] 错误，err:", err)
		return "", err
	}
	return addr, nil
}

// DelUserOnline 删除用户在线信息（存在即在线）
func DelUserOnline(userId int64) error {
	key := getUserKey(userId)
	_, err := public.Redis.Del(context.Background(), key).Result()
	if err != nil {
		zap.S().Error("[删除用户在线] 错误, err:", err)
		return err
	}
	return nil
}
