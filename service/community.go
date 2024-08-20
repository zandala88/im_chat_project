package service

import (
	"errors"
	"im/model"
	"log"
)

// 获取用户所加的群以及自己创建的群
func UserCommunities(userId int) ([]model.Community, error) {
	var communities []model.Community
	//var communityUsers []model.CommunityUsers
	result := DbEngine.Model(&model.CommunityUsers{}).Select("communities.*").Joins("left join communities on community_users.community_id = communities.id").Where("community_users.user_id = ?", userId).Scan(&communities)

	log.Printf("%v", communities)
	if result.Error != nil {
		return communities, result.Error
	}
	return communities, nil
}

// 用户加入群聊
func JoinCommunity(communityId, userId int) (bool, error) {
	// 查找用户是否已经添加了群聊
	communityUser := model.CommunityUsers{
		CommunityId: communityId,
		UserId:      userId,
	}

	// 查找群是否存在
	result := DbEngine.First(&model.Community{ID: communityId})
	if result.RowsAffected == 0 {
		// 未找到群聊
		return false, errors.New("未找到群组")
	}
	DbEngine.Where(&communityUser).Find(&communityUser)
	if communityUser.ID > 0 {
		return false, errors.New("您已经在群组里面")
	}

	result = DbEngine.Create(&communityUser)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// 创建群组
func CreateCommunity(ownerId int, name string) (model.Community, error) {
	community := model.Community{
		Name:    name,
		OwnerId: ownerId,
	}
	// 判断是否次用户已经创建过群组的功能, 如果有, 则不能创建
	DbEngine.Where("owner_id = ?", ownerId).First(&community)
	if community.ID > 0 {
		return community, errors.New("您已创建一个群组, 每个用户只能创建一个群组")
	}

	// 给用户添加群组, 并吧用户添加到群组里面
	DbEngine.Begin()

	communityResult := DbEngine.Create(&community)
	communityUser := model.CommunityUsers{
		CommunityId: community.ID,
		UserId:      ownerId,
	}

	communityUserResult := DbEngine.Create(&communityUser)
	DbEngine.Rollback()
	if communityResult.Error != nil || communityUserResult.Error != nil {
		DbEngine.Rollback()
		if communityResult.Error != nil {
			return community, communityResult.Error
		} else {
			return community, communityUserResult.Error
		}
	}
	DbEngine.Commit()
	return community, nil
}

// 查询群组的信息
func CommunityInfo(community *model.Community) (com *model.Community, err error) {
	// 查找
	DbEngine.First(&community)
	if community.ID == 0 {
		return community, errors.New("无效的参数")
	}

	return community, nil
}

// 获取用户的群组信息
func GetCommunitiesByUserId(userId int) ([]model.CommunityUsers, error) {
	var communityUsers []model.CommunityUsers
	result := DbEngine.Where("user_id = ?", userId).Find(&communityUsers)
	if result.Error != nil {
		return communityUsers, result.Error
	}

	return communityUsers, nil
}

// 获取群组的最新 10 条数据
func GetCommunityMessages(communityId int) []model.Message {
	var messages []model.Message
	DbEngine.Where("cmd = ?", model.CMD_GROUP).Where("to_id = ?", communityId).Order("id desc").Limit(10).Find(&messages)

	return messages
}
