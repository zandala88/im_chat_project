package service

import "im/server/cache"

func GetUserNextSeq(userId int64) (int64, error) {
	return cache.GetNextSeqId(cache.SeqObjectTypeUser, userId)
}

func GetUserNextSeqBatch(userIds []int64) ([]int64, error) {
	return cache.GetNextSeqIds(cache.SeqObjectTypeUser, userIds)
}
