package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"im/repo"
	"im/util"
)

type AddFriendRequest struct {
	FriendId int64 `json:"friendId"`
}

type AddFriendReply struct {
}

// AddFriend
// @Tags 好友
// @Summary 添加好友
// @accept json
// @Produce  json
// @Param token header string true "token"
// @Param friend body AddFriendRequest true "添加好友"
// @Success 200 {object} AddFriendReply
// @Router /add/friend [post]
func AddFriend(c *gin.Context) {
	req := &AddFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Errorf("[BindJSON ERROR] : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	userId := util.GetUid(c)
	if req.FriendId == userId {
		zap.S().Errorf("用户%d尝试添加自己为好友", userId)
		util.FailResp(c, "不能添加自己为好友")
		return
	}

	_, err := repo.GetUserByUserId(req.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Errorf("用户%d尝试添查询好友%d", userId, req.FriendId)
			util.FailResp(c, "该用户不存在")
			return
		}
		util.FailResp(c, "查询失败")
		return
	}

	_, err = repo.GetOneContact(userId, req.FriendId)
	if err == nil {
		zap.S().Errorf("用户%d已经添加过好友%d", userId, req.FriendId)
		util.FailResp(c, "已经添加过好友")
		return
	}

	err = repo.AddFriendCreate(userId, req.FriendId)
	if err != nil {
		zap.S().Errorf("repo.AddFriendCreate : %v", err)
		util.FailResp(c, "插入失败")
		return
	}
	util.SuccessResp(c, &AddFriendReply{})
}

type GetFriendRequest struct {
	FriendId int64 `json:"friendId"`
}

type GetFriendReply struct {
	UserName string `json:"username"`
	Mobile   string `json:"mobile"`
}

// GetFriend
// @Tags 好友
// @Summary 好友列表
// @accept json
// @Produce  json
// @Param friend body AddFriendRequest true "好友列表"
// @Param token header string true "token"
// @Success 200 {object} AddFriendReply
// @Router /get/friend [get]
func GetFriend(c *gin.Context) {
	req := &GetFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Errorf("[BindQuery ERROR] : %v", err)
		util.FailResp(c, err.Error())
		return
	}

	user, err := repo.GetUserByUserId(req.FriendId)
	if err != nil {
		zap.S().Errorf("repo.GetUserByUserId : %v", err)
		util.FailResp(c, "查询失败")
		return
	}

	util.SuccessResp(c, &GetFriendReply{
		UserName: user.UserName,
		Mobile:   user.Mobile,
	})
}

type GetFriendListRequest struct {
}

type GetFriendListReply struct {
	List []*GetFriendListSimple `json:"list"`
}

type GetFriendListSimple struct {
	Id     int64  `json:"id"`
	Avatar string `json:"avatar"`
}

// GetFriendList
// @Tags 好友
// @Summary 好友信息
// @accept json
// @Produce  json
// @Param friend body GetFriendListRequest true "好友列表"
// @Param token header string true "token"
// @Success 200 {object} GetFriendListReply
// @Router /get/friend/list [get]
func GetFriendList(c *gin.Context) {
	userId := util.GetUid(c)
	idList, err := repo.GetFriendList(userId)
	if err != nil {
		zap.S().Errorf("repo.GetFriendList : %v", err)
		util.FailResp(c, "查询失败")
		return
	}

	userList, err := repo.GetUserByUserIds(idList)
	if err != nil {
		zap.S().Errorf("repo.GetUserByUserIds : %v", err)
		util.FailResp(c, "查询失败")
		return
	}

	data := &GetFriendListReply{}
	for _, user := range userList {
		data.List = append(data.List, &GetFriendListSimple{
			Id: user.Id,
		})

	}

	util.SuccessResp(c, data)
}
