package service

import (
	"github.com/spf13/cast"
	"im/util"
)

func GetUserNextId(userId int64) (int64, error) {
	return util.UidGen.GetNextId(cast.ToString(userId))
}

// GetUserNextIdBatch 批量获取 seq
func GetUserNextIdBatch(userIds []int64) ([]int64, error) {
	businessIds := make([]string, len(userIds))
	for i, userId := range userIds {
		businessIds[i] = cast.ToString(userId)
	}
	return util.UidGen.GetNextIds(businessIds)
}
