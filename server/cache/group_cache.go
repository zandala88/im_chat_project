package cache

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"im/public"
	"time"
)

const (
	groupUserPrefix = "group_user_" // 群成员信息
	ttl2H           = 2 * 60 * 60   // 2h
)

func getGroupUserKey(groupId int64) string {
	return fmt.Sprintf("%s%d", groupUserPrefix, groupId)
}

// SetGroupUser 设置群成员
func SetGroupUser(groupId int64, userIds []int64) error {
	key := getGroupUserKey(groupId)
	values := make([]string, 0, len(userIds))
	for _, userId := range userIds {
		values = append(values, cast.ToString(userId))
	}
	_, err := public.Redis.SAdd(context.Background(), key, values).Result()
	if err != nil {
		zap.S().Error("[SetGroupUser] err = ", err)
		return err
	}
	_, err = public.Redis.Expire(context.Background(), key, ttl2H*time.Second).Result()
	if err != nil {
		zap.S().Error("[SetGroupUser] err = ", err)
		return err
	}
	return nil
}

// GetGroupUser 获取群成员
func GetGroupUser(groupId int64) ([]int64, error) {
	key := getGroupUserKey(groupId)
	result, err := public.Redis.SMembers(context.Background(), key).Result()
	if err != nil {
		zap.S().Error("[GetGroupUser] err = ", err)
		return nil, err
	}
	userIds := make([]int64, 0, len(result))
	for _, v := range result {
		userIds = append(userIds, cast.ToInt64(v))
	}
	return userIds, nil
}

// DeleteGroupUser 删除某个成员
func DeleteGroupUser(groupId, userId int64) error {
	key := getGroupUserKey(groupId)
	_, err := public.Redis.SRem(context.Background(), key, userId).Result()
	if err != nil {
		zap.S().Error("[DeleteGroupUser] err = ", err)
		return err
	}
	return nil
}

// DeleteGroupUserAll 删除所有成员
func DeleteGroupUserAll(groupId int64) error {
	key := getGroupUserKey(groupId)
	_, err := public.Redis.Del(context.Background(), key).Result()
	if err != nil {
		zap.S().Error("[DeleteGroupUserAll] err = ", err)
		return err
	}
	return nil
}
